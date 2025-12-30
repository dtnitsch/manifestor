package scanner

import (
	"fmt"
	"strings"

	"github.com/dtnitsch/manifestor/internal/manifest"
)

func AssertNoSkippedChildLeakage(m *manifest.Manifest) error {
	skippedDirs := map[string]struct{}{}

	for _, s := range m.Skipped {
		if s.IsDir {
			skippedDirs[s.Path] = struct{}{}
		}
	}

	for _, node := range m.Nodes {
		for dir := range skippedDirs {
			if node.Path != dir && strings.HasPrefix(node.Path, dir+"/") {
				return fmt.Errorf(
					"invariant violation: node %q appears under skipped dir %q",
					node.Path,
					dir,
				)
			}
		}
	}

	for _, skipped := range m.Skipped {
		for dir := range skippedDirs {
			if skipped.Path != dir && strings.HasPrefix(skipped.Path, dir+"/") {
				return fmt.Errorf(
					"invariant violation: skipped entry %q appears under skipped dir %q",
					skipped.Path,
					dir,
				)
			}
		}
	}

	return nil
}

