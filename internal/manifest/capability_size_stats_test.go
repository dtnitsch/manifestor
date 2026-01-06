package manifest

import (
	"testing"

)

func testManifestWithCapabilities(t *testing.T, sizeStats bool) *Manifest {
	t.Helper()

	return &Manifest{
		Manifest: ManifestMeta{
			Version: "0.2",
			Capabilities: Capabilities{
				Rollup: RollupCapabilities{
					SizeStats: sizeStats,
				},
			},
		},
		Nodes: []*Node{},
	}
}

func withSizeStats(r *Rollup, total, min, max, mean, median int64) {
	r.Size.Total = total
	r.Size.Min = min
	r.Size.Max = max
	r.Size.Mean = mean
	r.Size.Median = median
}


func TestCapabilitySizeStats_HappyPath(t *testing.T) {
	m := testManifestWithCapabilities(t, true)

    m.Nodes = []*Node{
        {
            Path:  "root",
            IsDir: true,
            Rollup: &Rollup{
                TotalFiles: 3,
            },
        },
    }
	withSizeStats(m.Nodes[0].Rollup, 100, 10, 60, 33, 30)

	violations, err := m.Validate(ValidateOptions{
		Strict: true,
	})

	if err != nil {
		t.Fatalf("unexpected fatal error: %v", err)
	}

	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d: %+v", len(violations), violations)
	}
}

