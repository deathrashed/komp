package engine

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/deathrashed/komp/internal/codec"
)

func List(archive string) ([]string, error) {
	c, ok := codec.ByExtension(archive)
	if !ok {
		return nil, fmt.Errorf("unknown archive type: %s", archive)
	}
	if _, err := exec.LookPath(c.Bin); err != nil {
		return nil, fmt.Errorf("%s not installed — brew install %s", c.Bin, c.BrewFormula)
	}
	args := substitute(c.ListArgs, map[string]string{"out": archive})
	out, err := exec.Command(c.Bin, args...).Output()
	if err != nil {
		return nil, fmt.Errorf("list %s: %w", archive, err)
	}
	var lines []string
	for _, l := range strings.Split(string(out), "\n") {
		l = strings.TrimSpace(l)
		if l != "" {
			lines = append(lines, l)
		}
	}
	return lines, nil
}
