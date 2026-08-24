package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "komp",
	Short:        "Unified archive & image toolkit",
	Long:         "komp — compress, add, list, extract, clean, convert archives; build dmg/iso/pkg images.",
	SilenceUsage: true,
	Args:         cobra.ArbitraryArgs,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(exitCodeOf(err))
	}
}

func exitCodeOf(err error) int {
	type coder interface{ ExitCode() int }
	if c, ok := err.(coder); ok {
		return c.ExitCode()
	}
	return 1
}
