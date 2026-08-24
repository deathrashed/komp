package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"komp/internal/recents"
	"komp/internal/xdg"
)

func TestLsRecentEmpty(t *testing.T) {
	if _, err := execLook("unzip"); err != nil {
		t.Skip("no unzip")
	}
	tmp := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmp)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	rootCmd.SetArgs([]string{"ls", "--recent"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for empty recents")
	}
}

func TestLsRecentOneEntry(t *testing.T) {
	if _, err := execLook("unzip"); err != nil {
		t.Skip("no unzip")
	}
	if _, err := execLook("zip"); err != nil {
		t.Skip("no zip")
	}

	tmp := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmp)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	src := filepath.Join(tmp, "a.txt")
	os.WriteFile(src, []byte("A"), 0o600)
	archive := filepath.Join(tmp, "a.zip")
	cmd := exec.Command("zip", archive, src)
	cmd.Dir = tmp
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("zip failed: %v\n%s", err, string(out))
	}

	store, _ := recents.Load(xdg.RecentsFile())
	store.Touch(archive)

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"ls", "--recent"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("err=%v out=%s", err, out.String())
	}
}
