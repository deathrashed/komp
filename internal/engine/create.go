package engine

import (
	"fmt"
	"os"
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
	if req.Separate {
		all := Result{}
		for _, in := range req.Inputs {
			r, err := createSingleGroup(c, Request{Inputs: []string{in}, Format: req.Format, OutputDir: req.OutputDir, Level: req.Level, DeleteOriginals: req.DeleteOriginals, DryRun: req.DryRun})
			if err != nil {
				return all, err
			}
			all.Outputs = append(all.Outputs, r.Outputs...)
			all.Plans = append(all.Plans, r.Plans...)
		}
		return all, nil
	}
	if c.Kind == codec.KindArchive {
		return createSingleGroup(c, req)
	}
	if len(req.Inputs) > 1 && c.Kind == codec.KindStream && !req.Each {
		return createWrapped(c, req)
	}
	return createDirect(c, req)
}

func createSingleGroup(c codec.Codec, req Request) (Result, error) {
	if len(req.Inputs) == 1 && c.Kind == codec.KindStream {
		return createDirect(c, req)
	}
	first := req.Inputs[0]
	ext := strings.TrimPrefix(c.Extensions[0], ".")
	out := OutputName(first, req.OutputDir, ext)
	if req.DryRun {
		return Result{Plans: []string{c.Bin + " \u2026 -> " + out}}, nil
	}
	stage := commonDir(req.Inputs)
	if c.Name == "tar" {
		args := []string{"-cf", out}
		for _, in := range req.Inputs {
			rel, _ := filepath.Rel(stage, in)
			args = append(args, rel)
		}
		cmd := exec.Command("tar", args...)
		cmd.Dir = stage
		if b, err := cmd.CombinedOutput(); err != nil {
			return Result{}, fmt.Errorf("tar create: %w: %s", err, b)
		}
		res := Result{Outputs: []string{out}}
		if err := finalize(req, req.Inputs, []string{out}); err != nil {
			return res, err
		}
		return res, nil
	}
	members := make([]string, 0, len(req.Inputs))
	for _, in := range req.Inputs {
		rel, _ := filepath.Rel(stage, in)
		members = append(members, rel)
	}
	m := map[string]string{"out": out, "level": levelToken(c, req.Level)}
	args := substitute(c.CreateArgs, m)
	args = append(args, members...)
	args = fixupArchiveArgs(c, args, out, members, stage)
	cmd := exec.Command(c.Bin, args...)
	cmd.Dir = stage
	if b, err := cmd.CombinedOutput(); err != nil {
		return Result{}, fmt.Errorf("create %s: %w: %s", out, err, b)
	}
	res := Result{Outputs: []string{out}}
	if err := finalize(req, req.Inputs, []string{out}); err != nil {
		return res, err
	}
	return res, nil
}

func fixupArchiveArgs(c codec.Codec, args []string, out string, members []string, stage string) []string {
	clean := make([]string, 0, len(args))
	for _, a := range args {
		switch a {
		case "{in}", "{indir}", "{inbase}":
			continue
		default:
			clean = append(clean, a)
		}
	}
	return clean
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
		if err := finalize(req, []string{in}, []string{out}); err != nil {
			return res, err
		}
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
	first := req.Inputs[0]
	ext := c.WrapExt
	if ext == "" {
		ext = ".tar" + c.Extensions[0]
	}
	out := OutputName(first, req.OutputDir, strings.TrimPrefix(ext, "."))
	tarPath := strings.TrimSuffix(out, c.Extensions[0])
	if !strings.HasSuffix(tarPath, ".tar") {
		tarPath += ".tar"
	}
	if tarPath == out {
		return Result{}, fmt.Errorf("unsupported wrap target %q", ext)
	}

	if req.DryRun {
		return Result{Plans: []string{"tar -cf " + tarPath + " … && " + c.Bin + " …"}}, nil
	}
	// 1. tar the inputs (common parent so members stay relative)
	root := commonDir(req.Inputs)
	tarArgs := []string{"-cf", tarPath}
	for _, in := range req.Inputs {
		rel, _ := filepath.Rel(root, in)
		tarArgs = append(tarArgs, rel)
	}
	cmd := exec.Command("tar", tarArgs...)
	cmd.Dir = root
	if b, err := cmd.CombinedOutput(); err != nil {
		return Result{}, fmt.Errorf("tar wrap: %w: %s", err, b)
	}
	// 2. stream-compress the tar (in place via codec), 3. drop intermediate tar when distinct
	defer func() { if tarPath != out { os.Remove(tarPath) } }()
	streamIn := tarPath
	// choose args: if codec expects {out} keep it, else strip it (brief strips unconditionally)
	var args []string
	hasOut := false
	for _, a := range c.CreateArgs {
		if a == "{out}" {
			hasOut = true
			break
		}
	}
	if hasOut {
		args = substitute(c.CreateArgs, map[string]string{"in": streamIn, "out": out, "level": levelToken(c, req.Level), "indir": filepath.Dir(streamIn), "inbase": filepath.Base(streamIn)})
	} else {
		args = substitute(stripOut(c.CreateArgs), map[string]string{"in": streamIn, "level": levelToken(c, req.Level)})
	}
	if b, err := exec.Command(c.Bin, args...).CombinedOutput(); err != nil {
		return Result{}, fmt.Errorf("compress wrap: %w: %s", err, b)
	}
	res := Result{Outputs: []string{out}}
	if err := finalize(req, []string{first}, []string{out}); err != nil {
		return res, err
	}
	return res, nil
}

func stripOut(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if a != "{out}" {
			out = append(out, a)
		}
	}
	return out
}

func commonDir(paths []string) string {
	if len(paths) == 1 {
		return filepath.Dir(paths[0])
	}
	dir := filepath.Dir(paths[0])
	for _, p := range paths[1:] {
		for !strings.HasPrefix(filepath.Dir(p)+"/", dir+"/") && dir != "/" && dir != "." {
			dir = filepath.Dir(dir)
		}
	}
	return dir
}

func finalize(req Request, inputs, outputs []string) error {
	if req.Backup {
		for _, o := range outputs {
			if exists(o) {
				if _, err := Backup(o); err != nil {
					return err
				}
			}
		}
	}
	if req.DeleteOriginals {
		for _, in := range inputs {
			os.Remove(in)
		}
	}
	return nil
}
