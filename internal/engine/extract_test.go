package engine

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestExtractZipRoundTrip(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.txt"); os.WriteFile(src, []byte("body"), 0o600)
	za := filepath.Join(dir, "in.zip")
	exec.Command("zip", "-j", za, src).Run()

	dest := filepath.Join(dir, "out")
	if err := Extract(za, dest, false); err != nil { t.Fatal(err) }
	b, err := os.ReadFile(filepath.Join(dest, "in.txt"))
	if err != nil || string(b) != "body" { t.Fatalf("round-trip broken: %v %q", err, b) }
}

func TestExtractNoClobber(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.txt"); os.WriteFile(src, []byte("new"), 0o600)
	za := filepath.Join(dir, "in.zip")
	exec.Command("zip", "-j", za, src).Run()

	dest := filepath.Join(dir, "out")
	os.MkdirAll(dest, 0o755)
	os.WriteFile(filepath.Join(dest, "in.txt"), []byte("old"), 0o600)
	err := Extract(za, dest, false)
	if err == nil { t.Fatal("must refuse without overwrite") }
	b, _ := os.ReadFile(filepath.Join(dest, "in.txt"))
	if string(b) != "old" { t.Fatal("existing file disturbed") }
	if err := Extract(za, dest, true); err != nil { t.Fatalf("overwrite pass: %v", err) }
}
