package scanner

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
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
        if err != nil {
            return fmt.Errorf("walk %q: %w", path, err)
        }

        // Context cancellation
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        // Apply filters
        if s.filters.Blocked(path, d) {
		    s.recordSkip(path, d, "blocked by filter", s.filters.MatchedRule(path, d))

            if d.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }

        info, err := d.Info()
        if err != nil {
            return fmt.Errorf("stat %q: %w", path, err)
        }

        node := &manifest.Node{
            Path:  path,
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

	// Adding skip details for output
	for _, s := range s.skipped {
		m.Skipped = append(m.Skipped, s)
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

