package engine

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"komp/internal/junk"
)

func TestCleanTarRebuilds(t *testing.T) {
	dir := t.TempDir()
	stage := filepath.Join(dir, "stage")
	os.MkdirAll(filepath.Join(stage, "keep"), 0o755)
	os.MkdirAll(filepath.Join(stage, ".git"), 0o755)
	os.WriteFile(filepath.Join(stage, "keep", "a.txt"), []byte("a"), 0o600)
	os.WriteFile(filepath.Join(stage, ".git", "HEAD"), []byte("ref"), 0o600)
	ta := filepath.Join(dir, "x.tar")
	c := exec.Command("tar", "-cf", ta, "-C", stage, "keep", ".git")
	if b, err := c.CombinedOutput(); err != nil { t.Fatalf("fixture: %v %s", err, b) }

	n, err := Clean(ta, []string{"vcs"})
	if err != nil { t.Fatal(err) }
	if n != 2 { t.Fatalf("want 2 removals, got %d", n) }
	lines, _ := List(ta)
	for _, l := range lines {
		if junk.Match("vcs", l) { t.Fatalf(".git survived: %s", l) }
	}
	if st, err := os.Stat(ta); err != nil || st.Size() == 0 { t.Fatal("tar broken") }
}
