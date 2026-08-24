package engine

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"komp/internal/codec"
)

func Add(archive string, files []string) error {
	c, ok := codec.ByExtension(archive)
	if !ok {
		return fmt.Errorf("unknown archive type: %s", archive)
	}
	if len(c.AddArgs) == 0 {
		return fmt.Errorf("%s cannot accept additions (single-stream compressor — recreate instead)", c.Name)
	}
	if _, err := exec.LookPath(c.Bin); err != nil {
		return fmt.Errorf("%s not installed — brew install %s", c.Bin, c.BrewFormula)
	}
	stage := ""
	if len(files) > 1 {
		stage = commonDir(files)
	} else {
		stage = filepath.Dir(files[0])
	}
	var members []string
	for _, f := range files {
		rel, _ := filepath.Rel(stage, f)
		members = append(members, rel)
	}
	args := substitute(c.AddArgs, map[string]string{"out": archive})
	for _, a := range args {
		_ = a
	}
	args = filterPlaceholders(args)
	args = append(args, members...)
	cmd := exec.Command(c.Bin, args...)
	cmd.Dir = stage
	if b, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("add to %s: %w: %s", archive, err, b)
	}
	return nil
}

func filterPlaceholders(args []string) []string {
	out := args[:0]
	for _, a := range args {
		switch a {
		case "{in}", "{indir}", "{inbase}", "{member}", "{dest}":
			continue
		default:
			out = append(out, a)
		}
	}
	return out
}
