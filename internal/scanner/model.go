package scanner

import "github.com/dtnitsch/manifestor/internal/filter"

type Scanner struct {
	opts    Options
	filters FilterSet
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
