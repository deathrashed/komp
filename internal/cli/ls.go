package cli

import (
	"github.com/spf13/cobra"
	"komp/internal/engine"
)

var lsCmd = &cobra.Command{
	Use:   "ls <archive>",
	Short: "List archive contents",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		lines, err := engine.List(args[0])
		if err != nil {
			return cliErr(err, classify(err))
		}
		for _, l := range lines {
			cmd.Println(l)
		}
		return nil
	},
}

func init() { rootCmd.AddCommand(lsCmd) }
