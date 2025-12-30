package output

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/dtnitsch/manifestor/internal/manifest"
)

func WriteJSON(path string, m *manifest.Manifest) error {
    f, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("create output file: %w", err)
    }
    defer f.Close()

    enc := json.NewEncoder(f)
    enc.SetIndent("", "  ")

    if err := enc.Encode(m); err != nil {
        return fmt.Errorf("encode manifest: %w", err)
    }
    return nil
}

