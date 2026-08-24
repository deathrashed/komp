package junk

import (
	"path"
	"strings"
)

var Groups = []string{"macos", "windows", "vcs", "hidden"}

type pred func(member string) bool

var preds = map[string][]pred{
	"macos": {
		func(m string) bool { return path.Base(m) == ".DS_Store" },
		func(m string) bool { return strings.HasPrefix(path.Base(m), "._") },
		func(m string) bool { return strings.HasPrefix(m, "__MACOSX/") || m == "__MACOSX" },
		func(m string) bool { return m == ".Spotlight-V100" || strings.HasPrefix(m, ".Spotlight-V100/") },
		func(m string) bool { return m == ".Trashes" || strings.HasPrefix(m, ".Trashes/") },
		func(m string) bool { return m == ".fseventsd" || strings.HasPrefix(m, ".fseventsd/") },
		func(m string) bool { return path.Base(m) == ".LSOverride" },
		func(m string) bool { return strings.HasPrefix(path.Base(m), "Icon\x0d") || path.Base(m) == "Icon?" },
	},
	"windows": {
		func(m string) bool { return path.Base(m) == "Thumbs.db" },
		func(m string) bool { return path.Base(m) == "desktop.ini" },
		func(m string) bool { return m == "$RECYCLE.BIN" || strings.HasPrefix(m, "$RECYCLE.BIN/") },
	},
	"vcs": {
		func(m string) bool {
			parts := strings.Split(m, "/")
			for _, p := range parts {
				if p == ".git" || p == ".svn" || p == ".hg" {
					return true
				}
			}
			return false
		},
	},
}

func Match(group, member string) bool {
	switch group {
	case "hidden":
		for g := range preds {
			if matchAny(g, member) { return false }
		}
		base := path.Base(member)
		return strings.HasPrefix(base, ".") && base != "." && base != ".."
	default:
		return matchAny(group, member)
	}
}

func matchAny(group, member string) bool {
	for _, p := range preds[group] {
		if p(member) { return true }
	}
	return false
}

func ValidGroups(requested []string) bool {
	for _, r := range requested {
		ok := r == "hidden"
		if !ok {
			for _, g := range Groups {
				if g == r { ok = true; break }
			}
		}
		if !ok { return false }
	}
	return true
}
