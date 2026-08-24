package selection

import (
	"fmt"
	"os/exec"
	"strings"
)

type Osascript struct{}

const script = `tell application "Finder"
set out to {}
repeat with i in selection
set end of out to POSIX path of (i as alias)
end repeat
set AppleScript's text item delimiters to linefeed
return out as text
end tell`

func (Osascript) Selection() ([]string, error) {
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return nil, fmt.Errorf("finder selection failed: %w", err)
	}
	s := strings.TrimRight(string(out), "\n")
	if s == "" {
		return nil, fmt.Errorf("nothing selected in Finder")
	}
	return strings.Split(s, "\n"), nil
}
