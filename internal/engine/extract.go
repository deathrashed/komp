package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	return fmt.Errorf("image extraction not yet implemented")
}
