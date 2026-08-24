package xdg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDir(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(base, "cfg"))
	got := ConfigDir()
	want := filepath.Join(base, "cfg", "komp")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestConfigDirDefault(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".config", "komp")
	if got := ConfigDir(); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
