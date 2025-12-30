package scanner

import (
	"os"
	"path/filepath"
	
	"github.com/dtnitsch/manifestor/internal/filter"
	"github.com/dtnitsch/manifestor/internal/manifest"
)

type Scanner struct {
	opts    Options
	filters FilterSet

	// Explicity store skipped folders for output
    skipped map[string]manifest.SkippedEntry
}

type Options struct {
	Root               string
	MaxWorkers         int
	FollowSymlinks     bool
	CollectInodes      bool
	CollectTimestamps  bool
	CollectFileCounts  bool
}

func New(opts Options, filters FilterSet) *Scanner {
	return &Scanner{
		opts:    opts,
		filters: filters,
	}
}

type FilterSet struct {
	Block []filter.Rule
	Allow []filter.Rule
}

func (f FilterSet) MatchedRule(path string, d os.DirEntry) *filter.Rule {
	base := filepath.Base(path)

	for _, r := range f.Block {
		if matchDir(r, path, base, d) {
			return &r
		}
	}
	return nil
}
