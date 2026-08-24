package image

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func BuildSfx(archive, configFile, module, out string) error {
	fout, err := os.Create(out)
	if err != nil {
		return err
	}
	defer fout.Close()
	for _, part := range []string{module, configFile, archive} {
		f, err := os.Open(part)
		if err != nil {
			return fmt.Errorf("sfx part %s: %w", part, err)
		}
		if _, err := io.Copy(fout, f); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return fout.Sync()
}

func ResolveModule(path string) (string, error) {
	if path != "" {
		return path, nil
	}
	home, _ := os.UserHomeDir()
	p := filepath.Join(home, ".config", "komp", "7zSD.sfx")
	if _, err := os.Stat(p); err == nil {
		return p, nil
	}
	return "", fmt.Errorf("no SFX module found (pass --sfx-module or place 7zSD.sfx in ~/.config/komp/)")
}

func substitute(args []string, m map[string]string) []string {
	out := make([]string, len(args))
	for i, a := range args {
		for k, v := range m {
			a = strings.ReplaceAll(a, "{"+k+"}", v)
		}
		out[i] = a
	}
	return out
}

func DefaultConfig(title string) string {
	return ";!@Install@!UTF-8!\r\nTitle=\"" + title + "\"\r\n;!@InstallEnd@!\r\n"
}
