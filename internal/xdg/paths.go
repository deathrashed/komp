package xdg

import (
	"os"
	"path/filepath"
)

func ConfigDir() string {
	if v := os.Getenv("XDG_CONFIG_HOME"); v != "" {
		return filepath.Join(v, "komp")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".komp"
	}
	return filepath.Join(home, ".config", "komp")
}

func RecentsFile() string { return filepath.Join(ConfigDir(), "recents.json") }
