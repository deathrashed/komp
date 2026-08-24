package ui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"komp/internal/codec"
	"komp/internal/recents"
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
	return runPicker(root)
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

func PickRecent(entries []recents.Entry) (string, error) {
	if !Interactive() {
		return "", errors.New("picking needs a terminal")
	}
	if len(entries) == 0 {
		return "", errors.New("no recent archives")
	}
	opts := make([]huh.Option[string], 0, len(entries))
	for _, e := range entries {
		label := e.Path + " — " + e.LastUsed.Format("2006-01-02 15:04")
		opts = append(opts, huh.NewOption(label, e.Path))
	}
	var choice string
	err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Recent archives").Options(opts...).Value(&choice),
	)).Run()
	return choice, err
}

func PickExistingArchive() (string, error) {
	if !Interactive() {
		return "", errors.New("file picking needs a terminal")
	}
	for {
		files, err := runPicker("")
		if err != nil {
			return "", err
		}
		var valid []string
		for _, f := range files {
			if _, ok := codec.ByExtension(f); ok {
				valid = append(valid, f)
			}
		}
		switch len(valid) {
		case 1:
			return valid[0], nil
		case 0:
			fmt.Fprintln(os.Stderr, "error: no valid archive selected — pick a zip, 7z, tar, etc.")
		default:
			fmt.Fprintln(os.Stderr, "error: select exactly one archive file")
		}
	}
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
