package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOutputNameBasic(t *testing.T) {
	got := OutputName("/x/notes.md", "/y/out", "zip")
	if got != "/y/out/notes.zip" {
		t.Fatalf("got %q", got)
	}
}

func TestOutputNameCollisionSuffix(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "notes.zip"), []byte("x"), 0o600)
	os.WriteFile(filepath.Join(dir, "notes-1.zip"), []byte("x"), 0o600)
	got := OutputName(filepath.Join(dir, "notes.md"), dir, "zip")
	if filepath.Base(got) != "notes-2.zip" {
		t.Fatalf("got %q", got)
	}
}

func TestOutputNameNoDoubleExt(t *testing.T) {
	// Stream compression appends the codec extension to the full filename.
	if got := OutputName("/x/a.md", "", "gz"); filepath.Base(got) != "a.md.gz" {
		t.Fatalf("stream ext should append: %q", got)
	}
}
