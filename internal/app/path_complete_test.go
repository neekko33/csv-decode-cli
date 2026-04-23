package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCollapseHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() error = %v", err)
	}

	got := collapseHome(filepath.Join(home, "Documents"), "~/Do")
	if got != "~/Documents" {
		t.Fatalf("collapseHome() = %q, want %q", got, "~/Documents")
	}
}

func TestCompletePathKeepsTildeForHomeInput(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() error = %v", err)
	}

	searchRoot, err := os.MkdirTemp(home, "csv-decode-complete-*")
	if err != nil {
		t.Fatalf("MkdirTemp() error = %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(searchRoot)
	})

	targetName := "alpha-file.csv"
	targetPath := filepath.Join(searchRoot, targetName)
	if err := os.WriteFile(targetPath, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	relDir := filepath.Base(searchRoot)
	input := filepath.Join("~", relDir, "alp")
	want := filepath.Join("~", relDir, targetName)

	got, ok := completePath(input)
	if !ok {
		t.Fatalf("completePath(%q) should change result", input)
	}
	if got != want {
		t.Fatalf("completePath(%q) = %q, want %q", input, got, want)
	}
}
