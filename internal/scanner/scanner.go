package scanner

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
	"sort"
    "syscall"
    "time"

    "github.com/dtnitsch/manifestor/internal/filter"
    "github.com/dtnitsch/manifestor/internal/manifest"
)

func (s *Scanner) Scan(ctx context.Context) (*manifest.Manifest, error) {
    m := &manifest.Manifest{
        Root:      s.opts.Root,
        Generated: time.Now().UTC(),
    }

	s.skipped = make(map[string]manifest.SkippedEntry)

    err := filepath.WalkDir(s.opts.Root, func(path string, d os.DirEntry, err error) error {
		norm := s.normalizePath(path)

        if err != nil {
            return fmt.Errorf("walk %q: %w", norm, err)
        }

        // Context cancellation
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        // Apply filters
        if s.filters.Blocked(norm, d) {
			s.recordSkip(norm, d, "blocked by filter", s.filters.MatchedRule(norm, d))

            if d.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }

        info, err := d.Info()
        if err != nil {
            return fmt.Errorf("stat %q: %w", norm, err)
        }

        node := &manifest.Node{
            Path:  norm,
            IsDir: d.IsDir(),
        }

        if s.opts.CollectTimestamps {
            node.ModTimeUnix = info.ModTime().Unix()
        }

        if s.opts.CollectInodes {
            if stat, ok := info.Sys().(*syscall.Stat_t); ok {
                node.Inode = stat.Ino
            }
        }

        if d.IsDir() && s.opts.CollectFileCounts {
            files, dirs := countDir(path)
            node.FileCount = files
            node.SubdirCount = dirs
        }

        m.Nodes = append(m.Nodes, node)
        return nil
    })

    if err != nil {
        return nil, err
    }

	// Cleanup of skipped things
	if len(s.skipped) > 0 {
		// Adding skip details for output
		for _, s := range s.skipped {
			m.Skipped = append(m.Skipped, s)
		}

		// Sort Skipped deterministically
		sort.Slice(m.Skipped, func(i, j int) bool {
			return m.Skipped[i].Path < m.Skipped[j].Path
		})
	}


    return m, nil
}

func (s *Scanner) recordSkip(path string, d os.DirEntry, reason string, rule *filter.Rule) {
	entry := manifest.SkippedEntry{
		Path:   path,
		IsDir:  d.IsDir(),
		Reason: reason,
	}

	if rule != nil {
		entry.Rule = string(rule.Type) + ":" + rule.Pattern
	}

	// De-dupe automatically
	s.skipped[path] = entry
}

func (s *Scanner) normalizePath(path string) string {
	// Convert to relative path
	rel, err := filepath.Rel(s.opts.Root, path)
	if err != nil {
		// Extremely defensive: fallback to original
		return filepath.Clean(path)
	}

	rel = filepath.Clean(rel)

	// filepath.Rel(".", ".") returns "."
	// filepath.Rel(".", "./foo") returns "foo"
	if rel == "" {
		return "."
	}

	return rel
}


