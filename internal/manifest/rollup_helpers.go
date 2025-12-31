package manifest

import (
	"os"
	"sort"
	"strings"
)

func sortNodesByDepthDesc(nodes []*Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return depth(nodes[i].Path) > depth(nodes[j].Path)
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func depth(path string) int {
	if path == "." {
		return 0
	}
	return strings.Count(path, string(os.PathSeparator)) + 1
}

func median(vals []int64) int64 {
	n := len(vals)
	if n == 0 {
		return 0
	}
	if n%2 == 1 {
		return vals[n/2]
	}
	return (vals[n/2-1] + vals[n/2]) / 2
}

