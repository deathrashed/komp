package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateStreamSingleFileGzip(t *testing.T) {
	if _, err := execLook("gzip"); err != nil {
		t.Skip("gzip unavailable")
	}
	dir := t.TempDir()
	src := filepath.Join(dir, "notes.txt")
	os.WriteFile(src, []byte("hello hello hello"), 0o600)

	req := Request{Inputs: []string{src}, Format: "gzip"}
	res, err := Create(req)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(res.Outputs[0]) != "notes.txt.gz" {
		t.Fatalf("out=%q", res.Outputs[0])
	}
	if _, err := os.Stat(res.Outputs[0]); err != nil {
		t.Fatal("no output")
	}
}
