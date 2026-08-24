package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"komp/internal/engine"
	"komp/internal/recents"
	"komp/internal/selection"
	"komp/internal/ui"
	"komp/internal/xdg"
)

var (
	addLast   bool
	addRecent bool
)

var addCmd = &cobra.Command{
	Use:   "add <archive> [files...]",
	Short: "Add files into an existing archive",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := recents.Load(xdg.RecentsFile())
		if err != nil {
			return cliErr(err, 1)
		}

		var archive string
		switch {
		case addLast:
			r := store.Recent(1)
			if len(r) == 0 {
				return cliErr(errors.New("no recent archives yet"), 1)
			}
			archive = r[0].Path
		case addRecent:
			if !ui.Interactive() {
				return cliErr(errors.New("--recent needs a terminal"), 1)
			}
			archive, err = ui.PickRecent(store.Recent(20))
			if err != nil {
				return cliErr(err, 1)
			}
		default:
			if len(args) == 0 {
				if !ui.Interactive() {
					return cliErr(errors.New("pass an archive path"), 1)
				}
				archive, err = ui.PickExistingArchive()
				if err != nil {
					return cliErr(err, 1)
				}
			} else {
				archive = args[0]
			}
		}

		var files []string
		if len(args) > 1 {
			files = args[1:]
		}
		if flagFinder {
			sel, err := selection.Osascript{}.Selection()
			if err != nil {
				return cliErr(err, 1)
			}
			files = append(files, sel...)
		}
		if len(files) == 0 && ui.Interactive() {
			files, err = ui.PickFiles("")
			if err != nil {
				return cliErr(err, 1)
			}
		}
		if len(files) == 0 {
			return cliErr(errors.New("no files to add"), 1)
		}

		if err := engine.Add(archive, files); err != nil {
			return cliErr(err, classify(err))
		}
		if err := store.Touch(archive); err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), "warn:", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "added %d file(s) to %s\n", len(files), archive)
		return nil
	},
}

func init() {
	addCmd.Flags().BoolVar(&flagFinder, "finder", false, "add current Finder selection")
	addCmd.Flags().BoolVar(&addLast, "last", false, "target most recent archive")
	addCmd.Flags().BoolVar(&addRecent, "recent", false, "choose from recent archives")
	addCmd.MarkFlagsMutuallyExclusive("last", "recent")
	rootCmd.AddCommand(addCmd)
}
