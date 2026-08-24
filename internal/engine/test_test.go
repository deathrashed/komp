package engine

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestVerifyGoodZipPasses(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.txt"); os.WriteFile(src, []byte("a"), 0o600)
	za := filepath.Join(dir, "a.zip")
	exec.Command("zip", "-j", za, src).Run()
	if err := Verify(za); err != nil { t.Fatal(err) }
}

func TestVerifyCorruptZipFails(t *testing.T) {
	dir := t.TempDir()
	za := filepath.Join(dir, "bad.zip")
	os.WriteFile(za, []byte("PK\x03\x04 garbage garbage"), 0o600)
	if err := Verify(za); err == nil { t.Fatal("corruption undetected") }
}

func TestCleanRefusesCorruptSource(t *testing.T) {
	dir := t.TempDir()
	za := filepath.Join(dir, "bad.zip")
	os.WriteFile(za, []byte("PK\x03\x04 garbage"), 0o600)
	if _, err := Clean(za, []string{"macos"}); err == nil {
		t.Fatal("clean must pre-flight verify source")
	}
}
