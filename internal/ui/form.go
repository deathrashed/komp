package ui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"komp/internal/codec"
	"komp/internal/junk"
	"komp/internal/recents"
)

var execLook = exec.LookPath

// ErrBack signals the user wants to return to the previous menu.
var ErrBack = errors.New("back")

func IsAbort(err error) bool {
	return errors.Is(err, huh.ErrUserAborted) || errors.Is(err, ErrBack)
}

func Interactive() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// Op is an interactive operation that operates on an existing archive.
type Op string

const (
	OpAdd     Op = "add"
	OpList    Op = "ls"
	OpExtract Op = "un"
	OpClean   Op = "clean"
	OpTest    Op = "t"
	OpConvert Op = "cv"
)

func (op Op) accepts(ext string) bool {
	name := strings.ToLower(ext)
	c, ok := codec.ByExtension("x" + name)
	if !ok {
		return isImageExt(name) && op == OpExtract
	}
	switch op {
	case OpAdd:
		return c.CanAdd()
	case OpList:
		return c.CanList()
	case OpExtract:
		return c.CanExtract() || isImageExt(name)
	case OpClean:
		return c.CanClean()
	case OpTest:
		return c.CanTest()
	case OpConvert:
		return c.CanExtract()
	}
	return false
}

func isImageExt(ext string) bool {
	switch ext {
	case ".dmg", ".sparseimage", ".sparsebundle", ".iso":
		return true
	}
	return false
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
		huh.NewOption("Add", "add"),
		huh.NewOption("Peek", "ls"),
		huh.NewOption("Extract", "un"),
		huh.NewOption("Clean", "clean"),
		huh.NewOption("Test", "t"),
		huh.NewOption("Convert", "cv"),
		huh.NewOption("Build", "img"),
	}
	var choice string
	err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Select Action").Description("Choose what komp should do. esc quits.").Options(opts...).Height(9).Value(&choice),
	)).WithShowHelp(true).Run()
	return choice, err
}

// archiveOptions lists archives in dir that the op can work with.
func archiveOptions(op Op, dir string) []huh.Option[string] {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var opts []huh.Option[string]
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if op.accepts(e.Name()) {
			opts = append(opts, huh.NewOption(e.Name(), filepath.Join(dir, e.Name())))
		}
	}
	sort.Slice(opts, func(i, j int) bool { return opts[i].Key < opts[j].Key })
	return opts
}

// PickArchiveFor shows capability-filtered archives in the current directory,
// with a manual-path entry at the bottom. Returns ErrBack on esc.
func PickArchiveFor(op Op) (string, error) {
	if !Interactive() {
		return "", errors.New("archive picking needs a terminal")
	}
	cwd, _ := os.Getwd()
	opts := archiveOptions(op, cwd)
	opts = append(opts, huh.NewOption("Type path manually", "__type__"))
	var choice string
	err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Select archive").Description("Only archives this action supports are listed. esc goes back.").Options(opts...).Value(&choice),
	)).WithShowHelp(true).Run()
	if err != nil {
		return "", ErrBack
	}
	if choice == "__type__" {
		var path string
		err := huh.NewForm(huh.NewGroup(
			huh.NewInput().Title("Archive path").Description("Full path to the archive. esc goes back.").Value(&path).Placeholder("/path/to/archive.zip"),
		)).WithShowHelp(true).Run()
		if err != nil {
			return "", ErrBack
		}
		return strings.TrimSpace(path), nil
	}
	return choice, nil
}

// PickCleanSettings is one form: archive + junk groups, previous answers stay visible.
func PickCleanSettings() (string, []string, error) {
	if !Interactive() {
		return "", nil, errors.New("interactive mode needs a terminal")
	}
	archive, groups := "", []string{}
	archOpts := archiveOptions(OpClean, mustCwd())
	archOpts = append(archOpts, huh.NewOption("Type path manually", "__type__"))
	groupOpts := make([]huh.Option[string], 0, len(junk.Groups))
	for _, g := range junk.Groups {
		groupOpts = append(groupOpts, huh.NewOption(g, g))
	}
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title("Select archive").Description("Only cleanable archives are listed. esc goes back.").Options(archOpts...).Value(&archive),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().Title("Junk groups to strip").Description("space toggles, enter submits, esc goes back.").Options(groupOpts...).Value(&groups),
		),
	).WithShowHelp(true).Run()
	if err != nil {
		return "", nil, ErrBack
	}
	if archive == "__type__" {
		var path string
		err := huh.NewForm(huh.NewGroup(
			huh.NewInput().Title("Archive path").Description("Full path to the archive. esc goes back.").Value(&path).Placeholder("/path/to/archive.zip"),
		)).WithShowHelp(true).Run()
		if err != nil {
			return "", nil, ErrBack
		}
		archive = strings.TrimSpace(path)
	}
	return archive, groups, nil
}

// PickConvertSettings is one form: archive + target format (source format excluded).
func PickConvertSettings() (string, string, error) {
	if !Interactive() {
		return "", "", errors.New("interactive mode needs a terminal")
	}
	archive, to := "", ""
	archOpts := archiveOptions(OpConvert, mustCwd())
	archOpts = append(archOpts, huh.NewOption("Type path manually", "__type__"))
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title("Select archive").Description("Only cleanable archives are listed. esc goes back.").Options(archOpts...).Value(&archive),
		),
	).WithShowHelp(true).Run()
	if err != nil {
		return "", "", ErrBack
	}
	if archive == "__type__" {
		var path string
		err := huh.NewForm(huh.NewGroup(
			huh.NewInput().Title("Archive path").Description("Full path to the archive. esc goes back.").Value(&path).Placeholder("/path/to/archive.tar.gz"),
		)).WithShowHelp(true).Run()
		if err != nil {
			return "", "", ErrBack
		}
		archive = strings.TrimSpace(path)
	}
	srcCodec, _ := codec.ByExtension(archive)
	fmtOpts := []huh.Option[string]{}
	for _, c := range codec.Table() {
		if srcCodec.Name == c.Name {
			continue
		}
		label := c.Name
		if !available(c.Bin) {
			label += "  (install: brew install " + c.BrewFormula + ")"
		}
		fmtOpts = append(fmtOpts, huh.NewOption(label, c.Name))
	}
	err = huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Convert to").Description("Target format. esc goes back.").Options(fmtOpts...).Value(&to),
	)).WithShowHelp(true).Run()
	if err != nil {
		return "", "", ErrBack
	}
	return archive, to, nil
}

// PickImageSettings is one form: source folder + image type + volume name.
func PickImageSettings() (src, kind, volname string, err error) {
	if !Interactive() {
		return "", "", "", errors.New("interactive mode needs a terminal")
	}
	kindOpts := []huh.Option[string]{
		huh.NewOption("dmg", "dmg"),
		huh.NewOption("sparsebundle", "sparsebundle"),
		huh.NewOption("sparseimage", "sparseimage"),
		huh.NewOption("iso", "iso"),
		huh.NewOption("pkg", "pkg"),
	}
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Source folder").Description("Folder to build the image from. esc goes back.").Value(&src).Placeholder("/path/to/folder").Validate(func(s string) error {
				if s == "" {
					return errors.New("folder path is required")
				}
				st, err := os.Stat(s)
				if err != nil {
					return fmt.Errorf("not accessible: %v", err)
				}
				if !st.IsDir() {
					return errors.New("path is not a folder")
				}
				return nil
			}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().Title("Image type").Description("esc goes back.").Options(kindOpts...).Value(&kind),
			huh.NewInput().Title("Volume name").Description("Shown when the image is mounted.").Value(&volname).Placeholder("(defaults to folder name)"),
		),
	).WithShowHelp(true).Run()
	if err != nil {
		return "", "", "", ErrBack
	}
	return src, kind, strings.TrimSpace(volname), nil
}

func mustCwd() string {
	d, _ := os.Getwd()
	return d
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
