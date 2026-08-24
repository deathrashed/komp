package engine

import (
	"path/filepath"
	"strings"
	"testing"

	"komp/internal/junk"
)

func TestCleanZipRemovesJunkInPlace(t *testing.T) {
	za := makeJunkyZip(t)
	before, _ := filepath.Abs(za)
	n, err := Clean(before, []string{"macos", "vcs"})
	if err != nil { t.Fatal(err) }
	if n < 3 { t.Fatalf("removed too few: %d", n) }
	lines, _ := List(before)
	for _, l := range lines {
		if junk.Match("macos", l) || junk.Match("vcs", l) {
			t.Fatalf("junk survived: %s", l)
		}
	}
	found := false
	for _, l := range lines { if strings.Contains(l, "Thumbs.db") { found = true } }
	if !found { t.Fatal("unselected group must survive") }
}
