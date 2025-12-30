package scanner

import (
    "path/filepath"
    "strings"

	"github.com/dtnitsch/manifestor/internal/filter"

)

func (f FilterSet) Blocked(path string, _ any) bool {
    for _, r := range f.Block {
        if match(r, path) {
            for _, a := range f.Allow {
                if match(a, path) {
                    return false
                }
            }
            return true
        }
    }
    return false
}

func match(r filter.Rule, path string) bool {
    switch r.Type {
    case filter.Basename:
        return filepath.Base(path) == r.Pattern
    case filter.Path:
        return strings.HasPrefix(path, r.Pattern)
    default:
        return false
    }
}

