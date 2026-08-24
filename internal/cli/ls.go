package cli

import (
	"errors"

	"github.com/spf13/cobra"
	"komp/internal/engine"
	"komp/internal/recents"
	"komp/internal/xdg"
)

var (
	lsRecent bool
)

var lsCmd = &cobra.Command{
	Use:   "ls <archive>",
	Short: "List archive contents",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var archive string
		if lsRecent {
			store, err := recents.Load(xdg.RecentsFile())
			if err != nil {
				return cliErr(err, 1)
			}
			r := store.Recent(1)
			if len(r) == 0 {
				return cliErr(errors.New("no recent archives yet"), 1)
			}
			archive = r[0].Path
		} else {
			if len(args) == 0 {
				return cliErr(errors.New("pass an archive path"), 1)
			}
			archive = args[0]
		}

		lines, err := engine.List(archive)
		if err != nil {
			return cliErr(err, classify(err))
		}
		for _, l := range lines {
			cmd.Println(l)
		}
		return nil
	},
}

func init() {
	lsCmd.Flags().BoolVar(&lsRecent, "recent", false, "list the most recent archive")
	rootCmd.AddCommand(lsCmd)
}
