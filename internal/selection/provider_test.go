package selection

import "testing"

type fake struct {
	items []string
	err   error
}

func (f fake) Selection() ([]string, error) { return f.items, f.err }

func TestInterfaceMock(t *testing.T) {
	var p Provider = fake{items: []string{"/tmp/a"}}
	items, err := p.Selection()
	if err != nil || len(items) != 1 {
		t.Fatal("mock broken")
	}
}
