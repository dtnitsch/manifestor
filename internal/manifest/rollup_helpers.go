package manifest

import (
	"math"
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

func computePercentiles(sorted []int64) *Percentiles {
	n := len(sorted)
	if n == 0 {
		return nil
	}

	at := func(p float64) int64 {
		rank := int(math.Ceil(p * float64(n)))
		if rank < 1 {
			rank = 1
		}
		if rank > n {
			rank = n
		}
		return sorted[rank-1]
	}

	return &Percentiles{
		P50: at(0.50),
		P90: at(0.90),
		P99: at(0.99),
	}
}

