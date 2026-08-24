package image

import (
	"fmt"
	"os/exec"
)

type Builder struct {
	Name        string
	Bin         string
	CreateArgs  []string
	Ext         string
	BrewFormula string
}

type Vars struct {
	Volname  string
	Format   string
	Size     string
	Fs       string
	ID       string
	Version  string
}

func Table() []Builder {
	return []Builder{
		{Name: "dmg", Bin: "hdiutil", Ext: ".dmg",
			CreateArgs: []string{"create", "-srcfolder", "{src}", "-volname", "{volname}", "-format", "{format}", "{out}"}},
		{Name: "sparsebundle", Bin: "hdiutil", Ext: ".sparsebundle",
			CreateArgs: []string{"create", "-type", "SPARSEBUNDLE", "-size", "{size}", "-fs", "{fs}", "-volname", "{volname}", "{out}"}},
		{Name: "sparseimage", Bin: "hdiutil", Ext: ".sparseimage",
			CreateArgs: []string{"create", "-type", "SPARSEIMAGE", "-size", "{size}", "-fs", "{fs}", "-volname", "{volname}", "{out}"}},
		{Name: "iso", Bin: "hdiutil", Ext: ".iso",
			CreateArgs: []string{"makehybrid", "-iso", "-joliet", "-default-volume-name", "{volname}", "-o", "{out}", "{src}"}},
		{Name: "pkg", Bin: "pkgbuild", Ext: ".pkg",
			CreateArgs: []string{"--root", "{src}", "--identifier", "{id}", "--version", "{version}", "{out}"}},
	}
}

func ByName(name string) (Builder, bool) {
	for _, b := range Table() {
		if b.Name == name {
			return b, true
		}
	}
	return Builder{}, false
}

func Build(name, src, out string, v Vars) error {
	b, ok := ByName(name)
	if !ok {
		return fmt.Errorf("unknown image type %q", name)
	}
	if _, err := exec.LookPath(b.Bin); err != nil {
		hint := b.BrewFormula
		if hint != "" {
			hint = " — brew install " + hint
		}
		return fmt.Errorf("%s not available%s", b.Bin, hint)
	}
	switch b.Name {
	case "dmg":
		if v.Format == "" {
			v.Format = "UDZO"
		}
	case "sparsebundle", "sparseimage":
		if v.Size == "" {
			return fmt.Errorf("--size required (e.g. 4g)")
		}
		if v.Fs == "" {
			v.Fs = "APFS"
		}
	case "pkg":
		if v.ID == "" {
			return fmt.Errorf("--id required for pkg (e.g. com.riley.tool)")
		}
		if v.Version == "" {
			v.Version = "1.0"
		}
	}
	args := substitute(b.CreateArgs, map[string]string{
		"src": src, "out": out, "volname": v.Volname, "format": v.Format,
		"size": v.Size, "fs": v.Fs, "id": v.ID, "version": v.Version,
	})
	cmd := exec.Command(b.Bin, args...)
	if out2, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s build: %w: %s", name, err, out2)
	}
	return nil
}
