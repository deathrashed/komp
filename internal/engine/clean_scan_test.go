package engine

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func makeJunkyZip(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("zip"); err != nil { t.Skip("no zip") }
	dir := t.TempDir()
	stage := filepath.Join(dir, "stage")
	files := map[string]string{
		"app/readme.txt":      "r",
		"app/.DS_Store":       "j",
		"__MACOSX/app/._meta": "j",
		"win/Thumbs.db":       "j",
		".git/config":         "j",
		".hiddenrc":           "h",
	}
	for rel, body := range files {
		full := filepath.Join(stage, rel)
		os.MkdirAll(filepath.Dir(full), 0o755)
		os.WriteFile(full, []byte(body), 0o600)
	}
	za := filepath.Join(dir, "junky.zip")
	cmd := exec.Command("zip", "-r", za, ".")
	cmd.Dir = stage
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("fixture zip failed: %v %s", err, b)
	}
	return za
}

func TestCleanScanCounts(t *testing.T) {
	za := makeJunkyZip(t)
	counts, err := CleanScan(za)
	if err != nil { t.Fatal(err) }
	if counts["macos"] < 2 || counts["windows"] < 1 || counts["vcs"] < 1 || counts["hidden"] < 1 {
		t.Fatalf("counts off: %+v", counts)
	}
}
