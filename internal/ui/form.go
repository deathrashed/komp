package ui

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"komp/internal/codec"
)

// execLook is indirection for exec.LookPath to allow testing.
var execLook = exec.LookPath

func Interactive() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func PickFiles(root string) ([]string, error) {
	if !Interactive() {
		return nil, errors.New("file picking needs a terminal — pass paths as arguments")
	}
	_ = root
	var raw string
	if err := huh.NewForm(
		huh.NewGroup(huh.NewText().
			Title("Files to compress (one path per line)").
			Value(&raw)),
	).Run(); err != nil {
		return nil, err
	}
	var out []string
	for _, l := range strings.Split(raw, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out, nil
}

func PickFormat() (string, error) {
	opts := []huh.Option[string]{}
	for _, c := range codec.Table() {
		label := c.Name
		if !available(c.Bin) {
			label += "  (install: brew install " + c.BrewFormula + ")"
		}
		o := huh.NewOption(label, c.Name)
		opts = append(opts, o)
	}
	var choice string
	err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Format").Options(opts...).Value(&choice),
	)).Run()
	return choice, err
}

func ConfirmDelete() bool {
	v := false
	_ = huh.NewConfirm().Title("Delete originals after compressing?").Value(&v).Run()
	return v
}

func available(bin string) bool {
	if bin == "" {
		return false
	}
	_, err := execLook(bin)
	return err == nil
}
