package engine

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"komp/internal/codec"
)

type Request struct {
	Inputs          []string
	Format          string
	OutputDir       string
	DeleteOriginals bool
	Separate        bool
	Each            bool
	Level           string // fast|normal|max
	Backup          bool
	DryRun          bool
}

type Result struct {
	Outputs []string
	Plans   []string // dry-run lines
}

func execLook(bin string) (string, error) { return exec.LookPath(bin) }

func levelToken(c codec.Codec, lvl string) string {
	switch lvl {
	case "fast":
		return c.Fast
	case "max":
		return c.Max
	default:
		return c.Normal
	}
}

func substitute(args []string, m map[string]string) []string {
	out := make([]string, len(args))
	for i, a := range args {
		for k, v := range m {
			a = strings.ReplaceAll(a, "{"+k+"}", v)
		}
		out[i] = a
	}
	return out
}

// Create compresses req.Inputs with the named format.
// Stream codec + single input → direct stream compression.
// Stream codec + multiple inputs → tar-wrap unless req.Each.
func Create(req Request) (Result, error) {
	c, ok := codec.ByName(req.Format)
	if !ok {
		return Result{}, fmt.Errorf("unknown format %q", req.Format)
	}
	if _, err := exec.LookPath(c.Bin); err != nil {
		return Result{}, fmt.Errorf("%s not installed — brew install %s", c.Bin, c.BrewFormula)
	}
	if len(req.Inputs) > 1 && c.Kind == codec.KindStream && !req.Each {
		return createWrapped(c, req)
	}
	return createDirect(c, req)
}

func createDirect(c codec.Codec, req Request) (Result, error) {
	res := Result{}
	for _, in := range req.Inputs {
		ext := pickExt(c, in)
		out := OutputName(in, req.OutputDir, ext)
		args := substitute(c.CreateArgs, map[string]string{
			"in": in, "out": out, "indir": filepath.Dir(in),
			"inbase": filepath.Base(in), "level": levelToken(c, req.Level),
		})
		if req.DryRun {
			res.Plans = append(res.Plans, c.Bin+" "+strings.Join(args, " "))
			continue
		}
		cmd := exec.Command(c.Bin, args...)
		if b, err := cmd.CombinedOutput(); err != nil {
			return res, fmt.Errorf("%s: %w: %s", in, err, b)
		}
		res.Outputs = append(res.Outputs, out)
	}
	return res, nil
}

func pickExt(c codec.Codec, in string) string {
	if c.Kind == codec.KindStream {
		return strings.TrimPrefix(c.Extensions[0], ".")
	}
	return strings.TrimPrefix(c.Extensions[0], ".")
}

func createWrapped(c codec.Codec, req Request) (Result, error) {
	// implemented Task 10
	return createDirect(c, req)
}
