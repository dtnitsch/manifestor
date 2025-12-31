package scanner_test

import (
	"context"
	"path/filepath"
	"os"
	"strings"
	"testing"

	"github.com/dtnitsch/manifestor/internal/filter"
	"github.com/dtnitsch/manifestor/internal/manifest"
	"github.com/dtnitsch/manifestor/internal/scanner"
)

func setup(t *testing.T) (*manifest.Manifest, *scanner.Scanner) {
	t.Helper()

	root := t.TempDir()

	// Create a fake .git directory with children
	err := os.MkdirAll(filepath.Join(root, ".git", "objects"), 0755)
	if err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	opts := scanner.Options{
		Root: root,
	}

	filters := scanner.FilterSet{
		Block: []filter.Rule{
			{
				Type:    filter.Basename,
				Pattern: ".git",
			},
		},
	}

	s := scanner.New(opts, filters)

	m, err := s.Scan(context.Background())
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	return m, s
}

func TestSkippedDirectorySuppressesChildren(t *testing.T) {
	m, _ := setup(t)

	// Assert .git is recorded as skipped
	foundGit := false
	for _, s := range m.Skipped {
		if s.Path == ".git" {
			foundGit = true
		}
		if strings.HasPrefix(s.Path, ".git/") {
			t.Fatalf("child skipped entry leaked: %q", s.Path)
		}
	}

	if !foundGit {
		t.Fatalf(".git was not recorded as skipped")
	}

	// Assert no nodes under .git
	for _, n := range m.Nodes {
		if strings.HasPrefix(n.Path, ".git/") {
			t.Fatalf("node leaked under skipped dir: %q", n.Path)
		}
	}
}

func TestSkipDirPreventsTraversal(t *testing.T) {
	m, _ := setup(t)

	for _, n := range m.Nodes {
		if strings.Contains(n.Path, "objects") {
			t.Fatalf("scanner descended into skipped directory: %q", n.Path)
		}
	}
}

func TestSkipDirNeverDescends(t *testing.T) {

	opts := scanner.Options{Root: "."}
	filters := scanner.FilterSet{
		Block: []filter.Rule{
			{Type: filter.Basename, Pattern: "internal/build"},
		},
	}

	s := scanner.New(opts, filters)

	m, err := s.Scan(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, n := range m.Nodes {
		if strings.Contains(n.Path, "/objects/") {
			t.Fatalf("descended into skipped directory: %q", n.Path)
		}
	}
}



