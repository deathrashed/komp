package recents

import (
	"os"
	"path/filepath"
	"testing"
)

func tempFile(t *testing.T) string {
	return filepath.Join(t.TempDir(), "recents.json")
}

func TestTouchAndRecent(t *testing.T) {
	p := tempFile(t)
	s, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Touch("/a/big.zip"); err != nil {
		t.Fatal(err)
	}
	if err := s.Touch("/a/big.zip"); err != nil {
		t.Fatal(err)
	}
	s.Touch("/b/x.7z")
	got := s.Recent(10)
	if len(got) != 2 || got[0].Path != "/b/x.7z" || got[1].Uses != 2 {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestCap20(t *testing.T) {
	s, _ := Load(tempFile(t))
	for i := 0; i < 25; i++ {
		s.Touch(string(rune('a'+i)) + ".zip")
	}
	if len(s.Recent(100)) != 20 {
		t.Fatal("cap not enforced")
	}
}

func TestCorruptFileStartsFresh(t *testing.T) {
	p := tempFile(t)
	os.WriteFile(p, []byte("not json"), 0o600)
	s, err := Load(p)
	if err != nil || s == nil {
		t.Fatalf("want fresh store, got err=%v", err)
	}
}
