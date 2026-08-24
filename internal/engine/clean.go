package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"komp/internal/codec"
	"komp/internal/junk"
)

// CleanScan returns how many members each junk group matches in archive.
func CleanScan(archive string) (map[string]int, error) {
	lines, err := List(archive)
	if err != nil { return nil, err }
	members := parseMembers(archive, lines)
	counts := map[string]int{}
	for _, m := range members {
		for _, g := range junk.Groups {
			if junk.Match(g, m) { counts[g]++ }
		}
	}
	return counts, nil
}

// parseMembers extracts clean relative paths from codec-specific ls output.
func parseMembers(archive string, lines []string) []string {
	var out []string
	lower := strings.ToLower(archive)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" { continue }
		switch {
		case strings.HasSuffix(lower, ".zip"):
			if isZipHeader(l) { continue }
			out = append(out, l)
		case strings.HasSuffix(lower, ".7z"):
			if isSevenZHeader(l) { continue }
			if idx := strings.LastIndex(l, " "); idx >= 0 && idx < len(l)-1 {
				out = append(out, strings.TrimSpace(l[idx+1:]))
			}
		case strings.HasSuffix(lower, ".tar"):
			out = append(out, l)
		default:
			out = append(out, l)
		}
	}
	return out
}

func isZipHeader(l string) bool {
	for _, p := range []string{"Archive contains", "Total bytes", "length", "----", "files"} {
		if strings.HasPrefix(l, p) { return true }
	}
	return false
}

func isSevenZHeader(l string) bool {
	for _, p := range []string{"7-Zip", "Listing", "--", "Date", "Path", "Size", "Copyright"} {
		if strings.HasPrefix(l, p) { return true }
	}
	return strings.TrimSpace(l) == ""
}

// Clean removes members matching any of groups from archive.
// zip/7z: in-place member deletion through atomic temp rewrite.
// tar*: full rebuild-and-swap.
func Clean(archive string, groups []string) (int, error) {
	if !junk.ValidGroups(groups) {
		return 0, fmt.Errorf("unknown junk group in %v (valid: macos windows vcs hidden)", groups)
	}
	c, ok := codec.ByExtension(archive)
	if !ok { return 0, fmt.Errorf("unknown archive type: %s", archive) }

	if err := Verify(archive); err != nil {
		return 0, fmt.Errorf("pre-flight failed: %w", err)
	}

	lines, err := List(archive)
	if err != nil { return 0, err }
	members := parseMembers(archive, lines)

	var doomed []string
	for _, m := range members {
		for _, g := range groups {
			if junk.Match(g, m) { doomed = append(doomed, m); break }
		}
	}
	if len(doomed) == 0 { return 0, nil }

	switch c.Name {
	case "zip":
		err = AtomicReplace(archive, func(tmp string) error {
			b, err := exec.Command("cp", archive, tmp).CombinedOutput()
			if err != nil { return fmt.Errorf("%w: %s", err, b) }
			zargs := []string{"-d", tmp}
			zargs = append(zargs, doomed...)
			if b, err := exec.Command("zip", zargs...).CombinedOutput(); err != nil {
				return fmt.Errorf("zip -d: %w: %s", err, b)
			}
			return nil
		})
	case "7z":
		err = AtomicReplace(archive, func(tmp string) error {
			if b, err := exec.Command("cp", archive, tmp).CombinedOutput(); err != nil {
				return fmt.Errorf("%w: %s", err, b)
			}
			dargs := []string{"d", tmp}
			dargs = append(dargs, doomed...)
			if b, err := exec.Command("7z", dargs...).CombinedOutput(); err != nil {
				return fmt.Errorf("7z d: %w: %s", err, b)
			}
			return nil
		})
	case "tar":
		err = cleanTar(archive, doomed)
	default:
		return 0, fmt.Errorf("%s cannot be cleaned (stream compressor)", c.Name)
	}
	if err != nil { return 0, err }
	return len(doomed), nil
}

func cleanTar(archive string, doomed []string) error {
	return AtomicReplace(archive, func(tmp string) error {
		extractDir, err := os.MkdirTemp("", "komp-clean-*")
		if err != nil { return err }
		defer os.RemoveAll(extractDir)
		if b, err := exec.Command("tar", "-xf", archive, "-C", extractDir).CombinedOutput(); err != nil {
			return fmt.Errorf("untar: %w: %s", err, b)
		}
		for _, d := range doomed {
			target := filepath.Join(extractDir, filepath.FromSlash(d))
			if err := os.RemoveAll(target); err != nil { return err }
		}
		args := append([]string{"-cf", tmp}, walkMembers(extractDir)...)
		cmd := exec.Command("tar", args...)
		cmd.Dir = extractDir
		if b, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("retar: %w: %s", err, b)
		}
		return nil
	})
}

func walkMembers(root string) []string {
	var out []string
	filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil || p == root { return nil }
		rel, _ := filepath.Rel(root, p)
		out = append(out, rel)
		return nil
	})
	return out
}
