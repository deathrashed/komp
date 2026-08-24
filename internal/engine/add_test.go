package engine

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestAddToZip(t *testing.T) {
	if _, err := exec.LookPath("zip"); err != nil {
		t.Skip("no zip")
	}
	dir := t.TempDir()
	extra := filepath.Join(dir, "extra.txt")
	os.WriteFile(extra, []byte("more"), 0o600)
	za := filepath.Join(dir, "a.zip")
	f := filepath.Join(dir, "first.txt")
	os.WriteFile(f, []byte("1"), 0o600)
	exec.Command("zip", "-j", za, f).Run()

	before, _ := os.Stat(za)
	if err := Add(za, []string{extra}); err != nil {
		t.Fatal(err)
	}
	after, _ := os.Stat(za)
	if after.Size() <= before.Size() {
		t.Fatal("archive did not grow")
	}
	lines, _ := List(za)
	ok := false
	for _, l := range lines {
		if filepath.Base(l) == "extra.txt" {
			ok = true
		}
	}
	if !ok {
		t.Fatal("extra.txt not added")
	}
}

func TestAddRejectsStreamFormat(t *testing.T) {
	dir := t.TempDir()
	gz := filepath.Join(dir, "a.txt.gz")
	os.WriteFile(gz, []byte("junk"), 0o600)
	err := Add(gz, []string{"whatever"})
	if err == nil {
		t.Fatal("stream add must fail")
	}
}
