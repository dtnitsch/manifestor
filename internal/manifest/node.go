package manifest

type Node struct {
	Path        string `json:"path" yaml:"path"`
	IsDir       bool   `json:"is_dir,omitempty" yaml:"is_dir,omitempty"`

	// Raw filesystem facts
	Inode       uint64 `json:"inode,omitempty" yaml:"inode,omitempty"`
	MtimeUnix int64  `json:"mtime_unix,omitempty" yaml:"mtime_unix,omitempty"`
	SizeBytes   int64  `json:"size_bytes,omitempty" yaml:"size_bytes,omitempty"`

	// Immediate directory stats (scanner)
	FileCount   int `json:"file_count,omitempty" yaml:"file_count,omitempty"`
	DirectSubdirCount int `json:"direct_subdir_count,omitempty" yaml:"direct_subdir_count,omitempty"`

	// Derived aggregate statistics
	Rollup *Rollup `json:"rollup,omitempty" yaml:"rollup,omitempty"`
}

