package manifest

type Node struct {
	Path        string `json:"path"`
	IsDir       bool   `json:"is_dir"`

	// Raw filesystem facts
	Inode       uint64 `json:"inode,omitempty"`
	MtimeUnix int64  `json:"mtime_unix,omitempty"`
	SizeBytes   int64  `json:"size_bytes,omitempty"`

	// Immediate directory stats (scanner)
	FileCount   int `json:"file_count,omitempty"`
	DirectSubdirCount int `json:"direct_subdir_count,omitempty"`

	// Derived aggregate statistics
	Rollup *Rollup `json:"rollup,omitempty"`
}

