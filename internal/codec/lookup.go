package codec

import (
	"path"
	"strings"
)

func ByName(name string) (Codec, bool) {
	low := strings.ToLower(name)
	for _, c := range Table() {
		if strings.ToLower(c.Name) == low {
			return c, true
		}
	}
	return Codec{}, false
}

func ByExtension(p string) (Codec, bool) {
	base := strings.ToLower(path.Base(p))
	var best Codec
	found := false
	bestLen := 0
	for _, c := range Table() {
		for _, ext := range c.Extensions {
			if len(ext) > 0 && strings.HasSuffix(base, ext) && len(ext) > bestLen {
				best, found = c, true
				bestLen = len(ext)
			}
		}
	}
	return best, found
}
