package engine

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"komp/internal/codec"
)

// Extract pulls archive into dest. overwrite=false refuses when any member
// already exists at target.
func Extract(archive, dest string, overwrite bool) error {
	ext := strings.ToLower(filepath.Ext(archive))
	if ext == ".dmg" || ext == ".sparseimage" || ext == ".sparsebundle" {
		return extractImage(archive, dest, overwrite)
	}

	c, ok := codec.ByExtension(archive)
	if !ok { return fmt.Errorf("unknown archive type: %s", archive) }
	if len(c.ExtractArgs) == 0 {
		return fmt.Errorf("%s cannot be extracted by komp yet", c.Name)
	}
	bin := c.ExtractBin
	if bin == "" { bin = c.Bin }
	if _, err := exec.LookPath(bin); err != nil {
		return fmt.Errorf("%s not installed — brew install %s", bin, c.BrewFormula)
	}
	if err := os.MkdirAll(dest, 0o755); err != nil { return err }
	if !overwrite {
		lines, err := List(archive)
		if err != nil { return err }
		for _, m := range parseMembers(archive, lines) {
			if _, err := os.Stat(filepath.Join(dest, filepath.FromSlash(strings.SplitN(m, " -> ", 2)[0]))); err == nil {
				return fmt.Errorf("%s exists — pass --overwrite to replace it", m)
			}
		}
	}
	args := substitute(c.ExtractArgs, map[string]string{"out": archive, "dest": dest})
	cmd := exec.Command(bin, args...)
	cmd.Dir = dest
	if b, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("extract %s: %w: %s", archive, err, b)
	}
	return nil
}

func extractImage(archive, dest string, overwrite bool) error {
	if err := os.MkdirAll(dest, 0o755); err != nil { return err }
	mp, err := attachImage(archive)
	if err != nil { return err }
	defer detachImage(mp)

	entries, err := os.ReadDir(mp)
	if err != nil { return err }
	for _, e := range entries {
		src := filepath.Join(mp, e.Name())
		dst := filepath.Join(dest, e.Name())
		if !overwrite {
			if _, err := os.Stat(dst); err == nil {
				return fmt.Errorf("%s exists — pass --overwrite", e.Name())
			}
		}
		if err := copyTree(src, dst); err != nil { return err }
	}
	return detachImage(mp)
}

func attachImage(archive string) (string, error) {
	mountPoint, err := os.MkdirTemp("", "komp-mount-*")
	if err != nil { return "", err }
	out, err := exec.Command("hdiutil", "attach", archive,
		"-readonly", "-nobrowse", "-mountpoint", mountPoint).CombinedOutput()
	if err != nil {
		os.RemoveAll(mountPoint)
		return "", fmt.Errorf("attach: %w: %s", err, out)
	}
	// hdiutil attach prints the mount point as the last tab-delimited field
	realMP := parseMountPoint(string(out))
	if realMP != "" && realMP != mountPoint {
		os.RemoveAll(mountPoint)
		mountPoint = realMP
	}
	return mountPoint, nil
}

func parseMountPoint(out string) string {
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Split(line, "\t")
		if len(fields) >= 2 {
			last := strings.TrimSpace(fields[len(fields)-1])
			if strings.HasPrefix(last, "/Volumes/") {
				return last
			}
		}
	}
	return ""
}

func detachImage(mp string) error {
	if err := exec.Command("hdiutil", "detach", mp).Run(); err != nil {
		return fmt.Errorf("detach %s: %w", mp, err)
	}
	return nil
}

func copyTree(src, dst string) error {
	si, err := os.Lstat(src)
	if err != nil { return err }
	if si.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(src)
		if err != nil { return err }
		return os.Symlink(target, dst)
	}
	if !si.IsDir() {
		in, err := os.Open(src)
		if err != nil { return err }
		defer in.Close()
		out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, si.Mode().Perm())
		if err != nil { return err }
		defer out.Close()
		_, err = io.Copy(out, in)
		if err != nil { return err }
		return os.Chtimes(dst, time.Time{}, si.ModTime())
	}
	if err := os.MkdirAll(dst, si.Mode().Perm()); err != nil { return err }
	entries, err := os.ReadDir(src)
	if err != nil { return err }
	for _, e := range entries {
		if err := copyTree(filepath.Join(src, e.Name()), filepath.Join(dst, e.Name())); err != nil {
			return err
		}
	}
	return nil
}
