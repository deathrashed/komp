package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestTCmdReportsIntegrity(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.txt"); os.WriteFile(src, []byte("a"), 0o600)
	za := filepath.Join(dir, "a.zip")
	exec.Command("zip", "-j", za, src).Run()

	var out bytes.Buffer
	rootCmd.SetOut(&out); rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"t", za})
	if err := rootCmd.Execute(); err != nil { t.Fatalf("err=%v out=%s", err, out.String()) }
	if !strings.Contains(out.String(), "OK") { t.Fatalf("expected OK, got %s", out.String()) }
}
