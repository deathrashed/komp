package junk

import "testing"

func TestMacosGroup(t *testing.T) {
	cases := map[string]bool{
		".DS_Store": true, "__MACOSX/x": true, "._resource": true,
		".Spotlight-V100": true, ".fseventsd/fsevents": true,
		"notes.txt": false, "sub/.DS_Store": true, "Icon\r": true,
	}
	for member, want := range cases {
		if got := Match("macos", member); got != want {
			t.Fatalf("macos %q: got %v want %v", member, got, want)
		}
	}
}

func TestWindowsGroup(t *testing.T) {
	cases := map[string]bool{
		"Thumbs.db": true, "desktop.ini": true, "$RECYCLE.BIN/file": true,
		"photo.jpg": false,
	}
	for m, want := range cases {
		if got := Match("windows", m); got != want { t.Fatalf("windows %q: %v", m, got) }
	}
}

func TestVcsGroup(t *testing.T) {
	if !Match("vcs", ".git/config") { t.Fatal(".git must match") }
	if !Match("vcs", "src/.svn/entries") { t.Fatal(".svn nested must match") }
	if Match("vcs", "gitfile") { t.Fatal("false positive") }
}

func TestHiddenGroup(t *testing.T) {
	if !Match("hidden", ".npmrc") { t.Fatal("dotfile must match hidden") }
	if Match("hidden", ".DS_Store") { t.Fatal("already covered by macos — hidden excludes other groups") }
}

func TestUnknownGroupErrors(t *testing.T) {
	if Match("nope", ".DS_Store") { t.Fatal("unknown group must never match") }
}
