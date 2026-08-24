package engine

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractDmgCopiesContents(t *testing.T) {
	if _, err := exec.LookPath("hdiutil"); err != nil { t.Skip("no hdiutil") }
	dir := t.TempDir()
	payload := filepath.Join(dir, "payload")
	os.MkdirAll(payload, 0o755)
	os.WriteFile(filepath.Join(payload, "tool.txt"), []byte("t"), 0o600)
	dmg := filepath.Join(dir, "tiny.dmg")
	if b, err := exec.Command("hdiutil", "create", "-srcfolder", payload,
		"-fs", "HFS+", "-volname", "TINY", "-size", "2m", dmg).CombinedOutput(); err != nil {
		t.Fatalf("fixture dmg: %v %s", err, b)
	}

	dest := filepath.Join(dir, "out")
	if err := Extract(dmg, dest, false); err != nil { t.Fatal(err) }
	if _, err := os.Stat(filepath.Join(dest, "tool.txt")); err != nil {
		t.Fatal("contents not copied")
	}
	out, _ := exec.Command("hdiutil", "info").Output()
	if hasMountedVolume(out, "TINY") { t.Fatal("image left mounted") }
}

func hasMountedVolume(hdiutilInfo []byte, vol string) bool {
	mp := "/Volumes/" + vol
	for _, line := range strings.Split(string(hdiutilInfo), "\n") {
		if strings.Contains(line, mp) { return true }
	}
	return false
}
