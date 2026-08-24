package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDryRunTouchesNothing(t *testing.T) {
	if _, err := execLook("gzip"); err != nil {
		t.Skip("no gzip")
	}
	dir := t.TempDir()
	src := filepath.Join(dir, "a.txt")
	os.WriteFile(src, []byte("A"), 0o600)
	res, err := Create(Request{Inputs: []string{src}, Format: "gzip", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Plans) == 0 {
		t.Fatal("expected plan lines")
	}
	if _, err := os.Stat(src + ".gz"); !os.IsNotExist(err) {
		t.Fatal("dry-run created a file!")
	}
	if _, err := os.Stat(filepath.Join(dir, "a.txt")); err != nil {
		t.Fatal("input disturbed")
	}
}
