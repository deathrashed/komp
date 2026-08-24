package engine

import (
	"io"
	"os"
	"path/filepath"
)

// AtomicReplace runs work(tmp) on a sibling temp file, then renames tmp over
// target only if work succeeded.
func AtomicReplace(target string, work func(tmp string) error) error {
	tmp, err := os.CreateTemp(filepath.Dir(target), "."+filepath.Base(target)+".tmp*")
	if err != nil { return err }
	tmp.Close()
	defer os.Remove(tmp.Name())
	if err := work(tmp.Name()); err != nil { return err }
	return os.Rename(tmp.Name(), target)
}

// Backup copies path to "<path>.bak" (overwriting any previous backup).
func Backup(path string) (string, error) {
	bak := path + ".bak"
	in, err := os.Open(path)
	if err != nil { return "", err }
	defer in.Close()
	out, err := os.Create(bak)
	if err != nil { return "", err }
	defer out.Close()
	_, err = io.Copy(out, in)
	return bak, err
}
