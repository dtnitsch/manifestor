package scanner

import (
	"os"
    "path/filepath"
    "strings"

	"github.com/dtnitsch/manifestor/internal/filter"

)

func (f FilterSet) Blocked(path string, d os.DirEntry) bool {
	base := filepath.Base(path)

    for _, r := range f.Block {
        if matchDir(r, path, base, d) {
            for _, a := range f.Allow {
                if matchDir(a, path, base, d) {
                    return false
                }
            }
            return true
        }
    }
    return false
}

func matchDir(r filter.Rule, path, base string, d os.DirEntry) bool {
    switch r.Type {
    case filter.Basename:
        return d.IsDir() && base == r.Pattern
    case filter.Path:
        return strings.HasPrefix(path, r.Pattern)
    default:
        return false
    }
}

