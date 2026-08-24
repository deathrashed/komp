package codec

import (
	"strings"
	"testing"
)

func TestTableIntegrity(t *testing.T) {
	names := map[string]bool{}
	for _, c := range Table() {
		if names[c.Name] {
			t.Fatalf("dup codec %s", c.Name)
		}
		names[c.Name] = true
		if c.Bin == "" {
			t.Fatalf("%s: empty bin", c.Name)
		}
		if len(c.Extensions) == 0 {
			t.Fatalf("%s: no extensions", c.Name)
		}
		for _, e := range c.Extensions {
			if !strings.HasPrefix(e, ".") {
				t.Fatalf("%s: ext %q lacks dot", c.Name, e)
			}
		}
		if c.Kind == KindArchive && c.CreateArgs == nil {
			t.Fatalf("%s: archive without CreateArgs", c.Name)
		}
	}
}
