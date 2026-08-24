package engine

import (
	"fmt"
	"os/exec"

	"github.com/deathrashed/komp/internal/codec"
)

func Verify(archive string) error {
	c, ok := codec.ByExtension(archive)
	if !ok { return fmt.Errorf("unknown archive type: %s", archive) }
	if len(c.TestArgs) == 0 { return fmt.Errorf("no verifier for %s", c.Name) }
	if _, err := exec.LookPath(c.Bin); err != nil {
		return fmt.Errorf("%s not installed — brew install %s", c.Bin, c.BrewFormula)
	}
	args := substitute(c.TestArgs, map[string]string{"out": archive})
	if b, err := exec.Command(c.Bin, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("verification failed: %w: %s", err, b)
	}
	return nil
}
