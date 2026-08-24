package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/deathrashed/komp/internal/engine"
	"github.com/deathrashed/komp/internal/recents"
	"github.com/deathrashed/komp/internal/ui"
	"github.com/deathrashed/komp/internal/xdg"
)

var (
	unOut       string
	unHere      bool
	unOverwrite bool
)

var unCmd = &cobra.Command{
	Use:   "un <archive>...",
	Short: "Extract archives (and disk images)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, _ := recents.Load(xdg.RecentsFile())
		start := time.Now()
		failed := 0
		for _, archive := range args {
			dest := unOut
			if unHere {
				dest, _ = os.Getwd()
			}
			if dest == "" {
				if ui.Interactive() {
					d, err := ui.PickDestination(defaultDestFor(archive))
					if err != nil { failed++; continue }
					dest = d
				} else {
					dest = defaultDestFor(archive)
				}
			}
			if err := engine.Extract(archive, dest, unOverwrite); err != nil {
				failed++
				cliReport(cmd, archive, err)
				continue
			}
			store.Touch(archive)
			cmd.Println(archive, "->", dest)
		}
		ui.NotifyIfSlow(start, 2*time.Second, "komp", "extraction finished")
		if failed > 0 { return &coded{errors.New("some extractions failed"), 3} }
		return nil
	},
}

func init() {
	unCmd.Flags().StringVarP(&unOut, "output", "o", "", "destination directory")
	unCmd.Flags().BoolVar(&unHere, "here", false, "extract into current directory")
	unCmd.Flags().BoolVar(&unOverwrite, "overwrite", false, "replace existing files")
	unCmd.MarkFlagsMutuallyExclusive("output", "here")
	rootCmd.AddCommand(unCmd)
}

func defaultDestFor(archive string) string {
	base := filepath.Base(archive)
	return "./" + strings.TrimSuffix(base, filepath.Ext(base))
}
