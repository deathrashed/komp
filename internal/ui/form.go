package ui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"komp/internal/codec"
	"komp/internal/junk"
	"komp/internal/recents"
)

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
		opts = append(opts, huh.NewOption(label, c.Name))
	}
	var choice string
	err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Format").Options(opts...).Value(&choice),
	)).WithShowHelp(true).Run()
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
	)).WithShowHelp(true).Run()
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

func PickGroups() ([]string, error) {
	if !Interactive() {
		return nil, errors.New("group picking needs a terminal")
	}
	opts := make([]huh.Option[string], 0, len(junk.Groups))
	for _, g := range junk.Groups {
		opts = append(opts, huh.NewOption(g, g))
	}
	var chosen []string
	err := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[string]().Title("Junk groups to strip").Options(opts...).Value(&chosen),
	)).WithShowHelp(true).Run()
	return chosen, err
}

func PickDestination(defaultVal string) (string, error) {
	if !Interactive() {
		return defaultVal, nil
	}
	var d string
	err := huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("Destination").Value(&d).Placeholder(defaultVal),
	)).WithShowHelp(true).Run()
	if err != nil {
		return "", err
	}
	if d == "" {
		return defaultVal, nil
	}
	return d, nil
}

func PickCommand() (string, error) {
	if !Interactive() {
		return "", errors.New("command picking needs a terminal")
	}
	opts := []huh.Option[string]{
		huh.NewOption("Compress", "compress"),
		huh.NewOption("Add to archive", "add"),
		huh.NewOption("Peek (list contents)", "ls"),
		huh.NewOption("Extract", "un"),
		huh.NewOption("Clean", "clean"),
		huh.NewOption("Test integrity", "t"),
		huh.NewOption("Convert", "cv"),
		huh.NewOption("Build disk image", "img"),
	}
	var choice string
	err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Select Action").Options(opts...).Value(&choice),
	)).WithShowHelp(true).Run()
	return choice, err
}

func PickArchive() (string, error) {
	if !Interactive() {
		return "", errors.New("archive picking needs a terminal")
	}
	archives := listArchives(".")
	var opts []huh.Option[string]
	for _, a := range archives {
		opts = append(opts, huh.NewOption(a, a))
	}
	opts = append(opts,
		huh.NewOption("Type path manually", "__type__"),
	)
	var choice string
	err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Select archive").Options(opts...).Value(&choice),
	)).WithShowHelp(true).Run()
	if err != nil {
		return "", err
	}
	if choice == "__type__" {
		var path string
		_ = huh.NewForm(huh.NewGroup(
			huh.NewInput().Title("Archive path").Value(&path).Placeholder("/path/to/archive.zip"),
		)).WithShowHelp(true).Run()
		return strings.TrimSpace(path), nil
	}
	return choice, nil
}

func PickImageSource() (string, error) {
	if !Interactive() {
		return "", errors.New("source picking needs a terminal")
	}
	for {
		files, err := runPicker("")
		if err != nil {
			return "", err
		}
		for _, f := range files {
			st, err := os.Stat(f)
			if err == nil && st.IsDir() {
				return f, nil
			}
		}
		fmt.Fprintln(os.Stderr, "error: select a folder (not a file)")
	}
}

func listArchives(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var archives []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if _, ok := codec.ByExtension(e.Name()); ok {
			archives = append(archives, e.Name())
		}
	}
	sort.Strings(archives)
	return archives
}

func archiveExts() []string {
	seen := map[string]bool{}
	var out []string
	for _, c := range codec.Table() {
		for _, e := range c.Extensions {
			if !seen[e] {
				seen[e] = true
				out = append(out, e)
			}
		}
	}
	return out
}
