package engine

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestListZip(t *testing.T) {
	if _, err := exec.LookPath("zip"); err != nil {
		t.Skip("no zip")
	}
	dir := t.TempDir()
	src := filepath.Join(dir, "hello.txt")
	os.WriteFile(src, []byte("hi"), 0o600)
	za := filepath.Join(dir, "h.zip")
	exec.Command("zip", "-j", za, src).Run()
	lines, err := List(za)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, l := range lines {
		if strings.Contains(l, "hello.txt") {
			found = true
		}
	}
	if !found {
		t.Fatalf("member missing: %v", lines)
	}
}
