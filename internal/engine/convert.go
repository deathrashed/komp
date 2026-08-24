package engine

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/deathrashed/komp/internal/codec"
)

type Options struct {
	DeleteOriginals bool
	OutputDir       string
	Prefix          string
	Suffix          string
	DryRun          bool
}

var ErrDryRunOnly = fmt.Errorf("dry run")

func Convert(archive, target string, opt Options) (string, error) {
	tc, ok := codec.ByName(target)
	if !ok {
		return "", fmt.Errorf("unknown target format %q", target)
	}
	if _, ok := codec.ByExtension(archive); !ok {
		return "", fmt.Errorf("unknown source archive type: %s", archive)
	}
	if err := Verify(archive); err != nil {
		return "", fmt.Errorf("source failed pre-flight: %w", err)
	}

	srcExt := strings.ToLower(filepath.Ext(archive))
	if inner, from, to := splitInnerTar(archive, srcExt, tc); inner != "" {
		return convertReuseInnerTar(archive, from, to, opt)
	}

	work, err := os.MkdirTemp("", "komp-cv-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(work)

	stage := filepath.Join(work, "x")
	if err := Extract(archive, stage, true); err != nil {
		return "", fmt.Errorf("stage extraction: %w", err)
	}
	members := walkRel(stage)
	outName := convertName(archive, tc, opt)
	outPath := joinOut(archive, opt.OutputDir, outName)
	if opt.DryRun {
		return outPath, ErrDryRunOnly
	}

	if err := AtomicReplace(outPath, func(tmp string) error {
		return createFromDir(tc, tmp, members, stage)
	}); err != nil {
		return "", err
	}
	if err := Verify(outPath); err != nil {
		os.Remove(outPath)
		return "", fmt.Errorf("post-conversion verification failed: %w", err)
	}
	if opt.DeleteOriginals {
		os.Remove(archive)
	}
	return outPath, nil
}

func splitInnerTar(archive, srcExt string, tc codec.Codec) (inner, from, to string) {
	const t = ".tar"
	if !strings.HasSuffix(strings.ToLower(archive), t+srcExt) && !strings.Contains(strings.ToLower(archive), t+".") {
		return "", "", ""
	}
	for _, te := range tc.Extensions {
		if strings.HasPrefix(te, t) {
			return "tar", srcExt, strings.TrimPrefix(te, t)
		}
	}
	return "", "", ""
}

func convertReuseInnerTar(archive, from, to string, opt Options) (string, error) {
	fromCodec, ok1 := codec.ByExtension("x" + from)
	toCodec, ok2 := codec.ByExtension("x" + to)
	if !ok1 || !ok2 {
		return "", fmt.Errorf("unsupported tar-stream pair %s->%s", from, to)
	}
	outName := swapCompressedExt(filepath.Base(archive), from, to)
	outPath := joinOut(archive, opt.OutputDir, outName)
	if opt.DryRun {
		return outPath, ErrDryRunOnly
	}
	return outPath, AtomicReplace(outPath, func(tmp string) error {
		return recompressStream(archive, tmp, fromCodec, toCodec)
	})
}

func recompressStream(src, dst string, dec, enc codec.Codec) error {
	dargs := substitute(decExtractArgs(dec), map[string]string{"in": src})
	eargs := substitute(stripOut(enc.CreateArgs), map[string]string{"level": levelToken(enc, "normal"), "in": "-"})
	decCmd := exec.Command(dargs[0], dargs[1:]...)
	encCmd := exec.Command(enc.Bin, eargs...)
	var w bytes.Buffer
	pr, err := decCmd.StdoutPipe()
	if err != nil {
		return err
	}
	encCmd.Stdin, encCmd.Stdout = pr, &w
	if err := encCmd.Start(); err != nil {
		return err
	}
	if out, err := decCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("decompress: %w: %s", err, out)
	}
	if err := encCmd.Wait(); err != nil {
		return fmt.Errorf("compress: %w", err)
	}
	return os.WriteFile(dst, w.Bytes(), 0o644)
}

func decExtractArgs(c codec.Codec) []string {
	switch c.Bin {
	case "gzip":
		return []string{"gzip", "-dc"}
	case "bzip2":
		return []string{"bzip2", "-dc"}
	case "xz":
		return []string{"xz", "-dc"}
	case "zstd":
		return []string{"zstd", "-dc"}
	case "lz4":
		return []string{"lz4", "-dc"}
	case "brotli":
		return []string{"brotli", "-dc"}
	case "lzip":
		return []string{"lzip", "-dc"}
	case "lzma":
		return []string{"lzma", "-dc"}
	case "lzo":
		return []string{"lzop", "-dc"}
	case "snappy":
		return []string{"snzip", "-dc", "-d"}
	case "lrzip":
		return []string{"lrzip", "-dc"}
	}
	return []string{c.Bin, "-d"}
}

func convertName(archive string, tc codec.Codec, opt Options) string {
	base := filepath.Base(archive)
	ext := tc.Extensions[0]
	clean := strings.TrimSuffix(base, filepath.Ext(base))
	for _, e := range tc.Extensions {
		if strings.HasSuffix(base, e) {
			clean = strings.TrimSuffix(base, e)
			break
		}
	}
	name := opt.Prefix + clean + opt.Suffix + ext
	return name
}

func joinOut(archive, outDir, name string) string {
	if outDir != "" {
		return filepath.Join(outDir, name)
	}
	return filepath.Join(filepath.Dir(archive), name)
}

func walkRel(root string) []string {
	var out []string
	filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil || p == root {
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		out = append(out, rel)
		return nil
	})
	return out
}

func createFromDir(tc codec.Codec, out string, members []string, root string) error {
	switch tc.Name {
	case "tar":
		args := []string{"-cf", out}
		for _, m := range members {
			args = append(args, m)
		}
		cmd := exec.Command("tar", args...)
		cmd.Dir = root
		if b, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("tar create: %w: %s", err, b)
		}
		return nil
	case "zip":
		args := []string{"-r", out}
		args = append(args, members...)
		cmd := exec.Command("zip", args...)
		cmd.Dir = root
		if b, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("zip create: %w: %s", err, b)
		}
		return nil
	case "7z":
		args := []string{"a", out}
		args = append(args, members...)
		cmd := exec.Command("7z", args...)
		cmd.Dir = root
		if b, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("7z create: %w: %s", err, b)
		}
		return nil
	}
	return fmt.Errorf("convert target %q not yet implemented", tc.Name)
}

func swapCompressedExt(name, from, to string) string {
	lower := strings.ToLower(name)
	fromLower := strings.ToLower(from)
	toLower := strings.ToLower(to)
	if strings.HasSuffix(lower, fromLower) {
		return name[:len(name)-len(fromLower)] + toLower
	}
	if strings.HasSuffix(lower, ".tar"+fromLower) {
		return name[:len(name)-len(".tar"+fromLower)] + ".tar" + toLower
	}
	return name + toLower
}
