package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestUnCmdIntoExplicitDest(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "f.txt"); os.WriteFile(src, []byte("x"), 0o600)
	za := filepath.Join(dir, "a.zip")
	exec.Command("zip", "-j", za, src).Run()
	dest := filepath.Join(dir, "here")

	var out bytes.Buffer
	rootCmd.SetOut(&out); rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"un", za, "-o", dest})
	if err := rootCmd.Execute(); err != nil { t.Fatalf("err=%v out=%s", err, out.String()) }
	if _, err := os.Stat(filepath.Join(dest, "f.txt")); err != nil { t.Fatal("not extracted") }
}
