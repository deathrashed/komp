package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"komp/internal/engine"
	"komp/internal/ui"
)

var (
	cleanGroups string
	cleanDry    bool
	cleanBackup bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean <archive>...",
	Short: "Strip junk members from archives",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		groups := []string{}
		if cleanGroups != "" {
			for _, g := range strings.Split(cleanGroups, ",") {
				groups = append(groups, strings.TrimSpace(g))
			}
		}
		if len(groups) == 0 {
			if !ui.Interactive() {
				return cliErr(errors.New("--groups required when not interactive"), 2)
			}
			sel, err := ui.PickGroups()
			if err != nil { return cliErr(err, 1) }
			groups = sel
		}
		failed := 0
		for _, archive := range args {
			counts, err := engine.CleanScan(archive)
			if err != nil { failed++; cliReport(cmd, archive, err); continue }
			total := 0
			for _, g := range groups { total += counts[g] }
			if total == 0 {
				cmd.Printf("%s: nothing to clean\n", archive)
				continue
			}
			if cleanDry {
				cmd.Printf("%s: would remove %d member(s): %v\n", archive, total, groups)
				continue
			}
			if cleanBackup { engine.Backup(archive) }
			n, err := engine.Clean(archive, groups)
			if err != nil { failed++; cliReport(cmd, archive, err); continue }
			cmd.Printf("%s: removed %d member(s)\n", archive, n)
		}
		ui.NotifyIfSlow(start, 2*time.Second, "komp", "clean finished")
		if failed > 0 { return &coded{fmt.Errorf("%d archive(s) failed", failed), 3} }
		return nil
	},
}

func init() {
	cleanCmd.Flags().StringVar(&cleanGroups, "groups", "", "comma list: macos,windows,vcs,hidden")
	cleanCmd.Flags().BoolVar(&cleanDry, "dry-run", false, "show what would be removed")
	cleanCmd.Flags().BoolVar(&cleanBackup, "backup", false, "save <archive>.bak first")
	rootCmd.AddCommand(cleanCmd)
}

func cliReport(cmd *cobra.Command, archive string, err error) {
	fmt.Fprintf(cmd.ErrOrStderr(), "%s: %v\n", archive, err)
}
