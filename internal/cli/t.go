package cli

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/deathrashed/komp/internal/engine"
)

var tCmd = &cobra.Command{
	Use:   "t <archive>...",
	Short: "Test archive integrity",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		failed := 0
		for _, a := range args {
			if err := engine.Verify(a); err != nil {
				failed++
				cliReport(cmd, a, err)
				continue
			}
			cmd.Println(a, "OK")
		}
		if failed > 0 { return &coded{errSilent, 4} }
		return nil
	},
}

func init() { rootCmd.AddCommand(tCmd) }

var errSilent = errors.New("")
