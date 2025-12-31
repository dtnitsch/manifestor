package scanner

import (
	"path/filepath"
	"os"
)


func countDir(path string, shouldInclude func(string) bool) (files, dirs int) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return 0, 0
    }

	for _, e := range entries {
	    full := filepath.Join(path, e.Name())

	    if !shouldInclude(full) {
	        continue
	    }

	    if e.IsDir() {
	        dirs++
	    } else {
	        files++
	    }
	}
    return
}

