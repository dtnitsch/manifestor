package manifest

type Node struct {
    Path         string `json:"path"`
    IsDir        bool   `json:"is_dir"`

    Inode        uint64 `json:"inode,omitempty"`
    ModTimeUnix  int64  `json:"mtime_unix,omitempty"`

    FileCount    int    `json:"file_count,omitempty"`
    SubdirCount  int    `json:"subdir_count,omitempty"`
}

