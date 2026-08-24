package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateCmdSilentZip(t *testing.T) {
	if _, err := execLook("zip"); err != nil {
		t.Skip("no zip")
	}
	dir := t.TempDir()
	src := filepath.Join(dir, "a.txt")
	os.WriteFile(src, []byte("A"), 0o600)
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"-f", "zip", "-o", dir, src})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("err=%v out=%s", err, out.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "a.zip")); err != nil {
		t.Fatal("no archive")
	}
}
