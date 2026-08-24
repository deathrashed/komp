package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"komp/internal/codec"
	"komp/internal/engine"
	"komp/internal/selection"
	"komp/internal/ui"
)

var (
	flagFormat string
	flagDelete bool
	flagOut    string
	flagLevel  string
	flagSep    bool
	flagEach   bool
	flagFinder bool
	flagDryRun bool
	flagBackup bool
)

var createCmd = &cobra.Command{
	Use:   "komp [format] [files...]",
	Short: "Compress files/folders",
	RunE:  runCreate,
}

func init() {
	createCmd.Flags().StringVarP(&flagFormat, "format", "f", "", "codec name (also: first bare arg)")
	createCmd.Flags().BoolVarP(&flagDelete, "delete", "d", false, "remove originals after success")
	createCmd.Flags().StringVarP(&flagOut, "output", "o", "", "output directory")
	createCmd.Flags().StringVarP(&flagLevel, "level", "L", "normal", "fast|normal|max")
	createCmd.Flags().BoolVar(&flagSep, "separate", false, "one archive per input")
	createCmd.Flags().BoolVar(&flagEach, "each", false, "streams: compress each input separately")
	createCmd.Flags().BoolVar(&flagFinder, "finder", false, "use current Finder selection as input")
	createCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "print plan, touch nothing")
	createCmd.Flags().BoolVar(&flagBackup, "backup", false, "back up overwritten targets")
	rootCmd.AddCommand(createCmd)

	// Also expose flags on root so `komp -f zip file` works (non-TTY piped path).
	// Shares same backing vars so behavior is identical.
	rootCmd.Flags().StringVarP(&flagFormat, "format", "f", "", "codec name (also: first bare arg)")
	rootCmd.Flags().BoolVarP(&flagDelete, "delete", "d", false, "remove originals after success")
	rootCmd.Flags().StringVarP(&flagOut, "output", "o", "", "output directory")
	rootCmd.Flags().StringVarP(&flagLevel, "level", "L", "normal", "fast|normal|max")
	rootCmd.Flags().BoolVar(&flagSep, "separate", false, "one archive per input")
	rootCmd.Flags().BoolVar(&flagEach, "each", false, "streams: compress each input separately")
	rootCmd.Flags().BoolVar(&flagFinder, "finder", false, "use current Finder selection as input")
	rootCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "print plan, touch nothing")
	rootCmd.Flags().BoolVar(&flagBackup, "backup", false, "back up overwritten targets")
	rootCmd.RunE = runRoot
}

func runRoot(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && ui.Interactive() {
		choice, err := ui.PickCommand()
		if err != nil {
			return cliErr(err, 1)
		}
		switch choice {
		case "compress":
			return runCreate(cmd, args)
		case "add":
			archive, err := ui.PickExistingArchive()
			if err != nil {
				return cliErr(err, 1)
			}
			return addCmd.RunE(addCmd, []string{archive})
		case "ls":
			archive, err := ui.PickArchive()
			if err != nil {
				return cliErr(err, 1)
			}
			return lsCmd.RunE(lsCmd, []string{archive})
		case "un":
			archive, err := ui.PickArchive()
			if err != nil {
				return cliErr(err, 1)
			}
			return unCmd.RunE(unCmd, []string{archive})
		case "clean":
			archive, err := ui.PickArchive()
			if err != nil {
				return cliErr(err, 1)
			}
			return cleanCmd.RunE(cleanCmd, []string{archive})
		case "t":
			archive, err := ui.PickArchive()
			if err != nil {
				return cliErr(err, 1)
			}
			return tCmd.RunE(tCmd, []string{archive})
		case "cv":
			archive, err := ui.PickArchive()
			if err != nil {
				return cliErr(err, 1)
			}
			return cvCmd.RunE(cvCmd, []string{archive})
		case "img":
			src, err := ui.PickFiles("")
			if err != nil {
				return cliErr(err, 1)
			}
			if len(src) == 0 {
				return cliErr(errors.New("no source folder selected"), 1)
			}
			return imgCmd.RunE(imgCmd, []string{src[0]})
		}
	}
	return runCreate(cmd, args)
}

func runCreate(cmd *cobra.Command, args []string) error {
	format := flagFormat
	var files []string
	for _, a := range args {
		if _, known := codec.ByName(a); known && format == "" {
			format = a // preset: bare positional codec
		} else {
			files = append(files, a)
		}
	}
	if flagFinder {
		sel, err := selection.Osascript{}.Selection()
		if err != nil {
			return cliErr(err, 1)
		}
		files = append(files, sel...)
	}
	if len(files) == 0 && ui.Interactive() {
		picked, err := ui.PickFiles("")
		if err != nil {
			return cliErr(err, 1)
		}
		files = picked
	}
	if len(files) == 0 {
		return cliErr(errors.New("no input files (empty Finder selection?)"), 1)
	}
	if format == "" {
		if !ui.Interactive() {
			return cliErr(errors.New("no format given — pass -f or a bare codec name"), 2)
		}
		f, err := ui.PickFormat()
		if err != nil {
			return cliErr(err, 1)
		}
		format = f
	}
	if !flagDelete && ui.Interactive() && flagFormat == "" {
		flagDelete = ui.ConfirmDelete() // only prompt when format wasn't preset
	}
	res, err := engine.Create(engine.Request{
		Inputs: files, Format: format, OutputDir: flagOut, DeleteOriginals: flagDelete,
		Separate: flagSep, Each: flagEach, Level: flagLevel, Backup: flagBackup, DryRun: flagDryRun,
	})
	if err != nil {
		return cliErr(err, classify(err))
	}
	printResult(cmd, res)
	return nil
}

func cliErr(err error, code int) error { return &coded{err, code} }

type coded struct {
	error
	code int
}

func (c *coded) ExitCode() int { return c.code }

func classify(err error) int {
	msg := err.Error()
	switch {
	case contains(msg, "not installed"):
		return 2
	case contains(msg, "unknown format"):
		return 2
	default:
		return 3
	}
}

func contains(s, sub string) bool { return len(s) >= len(sub) && indexOf(s, sub) >= 0 }
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func printResult(cmd *cobra.Command, res engine.Result) {
	for _, p := range res.Plans {
		fmt.Fprintln(os.Stderr, "plan:", p)
	}
	for _, o := range res.Outputs {
		fmt.Fprintln(cmd.OutOrStdout(), o)
	}
}

func execLook(bin string) (string, error) { return exec.LookPath(bin) }
