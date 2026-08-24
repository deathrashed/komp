package engine

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"komp/internal/codec"
)

// OutputName builds the output path in outDir (or the input's dir when
// empty). Archive formats replace the input's extension ("<name>.<ext>");
// stream formats append to the full filename ("<name>.<in-ext>.<ext>", so
// notes.md -> notes.md.gz). Collisions get -1, -2… inserted before ".<ext>".
func OutputName(input, outDir, ext string) string {
	dir := outDir
	if dir == "" {
		dir = filepath.Dir(input)
	}
	base := filepath.Base(input)
	if kindIsArchive(ext) {
		base = strings.TrimSuffix(base, filepath.Ext(base))
	}
	name := base + "." + ext
	full := filepath.Join(dir, name)
	for i := 1; exists(full); i++ {
		full = filepath.Join(dir, base+"-"+strconv.Itoa(i)+"."+ext)
	}
	return full
}

// kindIsArchive reports whether ext names a container codec. Unknown
// extensions (e.g. a synthesized "tar.gz") are treated as archives, which
// only matters when the input itself has an extension.
func kindIsArchive(ext string) bool {
	if c, ok := codec.ByName(ext); ok {
		return c.Kind == codec.KindArchive
	}
	for _, c := range codec.Table() {
		for _, e := range c.Extensions {
			if strings.TrimPrefix(e, ".") == ext && c.Kind != codec.KindArchive {
				return false
			}
		}
	}
	return true
}

func exists(p string) bool { _, err := os.Stat(p); return err == nil }
