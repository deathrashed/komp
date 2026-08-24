package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicSwapReplacesOnlyOnSuccess(t *testing.T) {
	target := filepath.Join(t.TempDir(), "a.zip")
	os.WriteFile(target, []byte("old"), 0o600)
	err := AtomicReplace(target, func(tmp string) error {
		return os.WriteFile(tmp, []byte("new"), 0o600)
	})
	if err != nil { t.Fatal(err) }
	b, _ := os.ReadFile(target)
	if string(b) != "new" { t.Fatal("swap did not happen") }
}

func TestAtomicSwapKeepsOldOnFailure(t *testing.T) {
	target := filepath.Join(t.TempDir(), "a.zip")
	os.WriteFile(target, []byte("old"), 0o600)
	_ = AtomicReplace(target, func(tmp string) error { return os.ErrPermission })
	b, _ := os.ReadFile(target)
	if string(b) != "old" { t.Fatal("original must survive failure") }
}

func TestBackupCopies(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.zip")
	os.WriteFile(src, []byte("data"), 0o600)
	bak, err := Backup(src)
	if err != nil { t.Fatal(err) }
	if filepath.Base(bak) != "a.zip.bak" { t.Fatalf("bak=%q", bak) }
}
