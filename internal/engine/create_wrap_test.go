package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateMultiInputStreamWrapsInTar(t *testing.T) {
	if _, err := execLook("tar"); err != nil { t.Skip("no tar") }
	if _, err := execLook("gzip"); err != nil { t.Skip("no gzip") }
	dir := t.TempDir()
	a := filepath.Join(dir, "a.txt"); b := filepath.Join(dir, "b.txt")
	os.WriteFile(a, []byte("A"), 0o600); os.WriteFile(b, []byte("B"), 0o600)
	req := Request{Inputs: []string{a, b}, Format: "gzip", OutputDir: dir}
	res, err := Create(req)
	if err != nil { t.Fatal(err) }
	if filepath.Base(res.Outputs[0]) != "a.tar.gz" { t.Fatalf("wrap name=%q", res.Outputs[0]) }
}

func TestCreateEachSplits(t *testing.T) {
	if _, err := execLook("gzip"); err != nil { t.Skip("no gzip") }
	dir := t.TempDir()
	a := filepath.Join(dir, "a.txt"); b := filepath.Join(dir, "b.txt")
	os.WriteFile(a, []byte("A"), 0o600); os.WriteFile(b, []byte("B"), 0o600)
	req := Request{Inputs: []string{a, b}, Format: "gzip", Each: true, OutputDir: dir}
	res, err := Create(req)
	if err != nil { t.Fatal(err) }
	if len(res.Outputs) != 2 { t.Fatalf("want 2 outputs, got %v", res.Outputs) }
}
