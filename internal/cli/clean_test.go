package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"komp/internal/engine"
)

func TestCleanCmdPipedWithGroups(t *testing.T) {
	dir := t.TempDir()
	za := filepath.Join(dir, "a.zip")
	f := filepath.Join(dir, "f.txt"); os.WriteFile(f, []byte("f"), 0o600)
	j := filepath.Join(dir, ".DS_Store"); os.WriteFile(j, []byte("j"), 0o600)
	exec.Command("zip", "-j", za, f).Run()
	exec.Command("zip", "-j", za, j).Run()

	var out bytes.Buffer
	rootCmd.SetOut(&out); rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"clean", "--groups", "macos", za})
	if err := rootCmd.Execute(); err != nil { t.Fatalf("err=%v out=%s", err, out.String()) }
	lines, _ := engine.List(za)
	for _, l := range lines {
		if strings.Contains(l, ".DS_Store") { t.Fatal("junk survived") }
	}
}
