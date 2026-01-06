package output

import (
	"fmt"
	"os"

	"github.com/dtnitsch/manifestor/internal/manifest"
	"gopkg.in/yaml.v3"
)

func WriteYAML(path string, m *manifest.Manifest) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)

	if err := enc.Encode(m); err != nil {
		return fmt.Errorf("encode manifest: %w", err)
	}
	return nil
}
