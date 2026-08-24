package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/deathrashed/komp/internal/engine"
	"github.com/deathrashed/komp/internal/ui"
)

type cvFlags struct {
	To             string
	Prefix         string
	Suffix         string
	OutDir         string
	DeleteOriginal bool
	DryRun         bool
}

var cvOpt cvFlags

var cvCmd = &cobra.Command{
	Use:   "cv <archive>",
	Short: "Recompress an archive to another format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		archive := args[0]
		if cvOpt.To == "" {
			pick, err := ui.PickFormat()
			if err != nil {
				return err
			}
			cvOpt.To = pick
		}
		out, err := engine.Convert(archive, cvOpt.To, engine.Options{
			Prefix:          cvOpt.Prefix,
			Suffix:          cvOpt.Suffix,
			OutputDir:       cvOpt.OutDir,
			DeleteOriginals: cvOpt.DeleteOriginal,
			DryRun:          cvOpt.DryRun,
		})
		if err == engine.ErrDryRunOnly {
			fmt.Println(out)
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println(out)
		return nil
	},
}

func init() {
	cvCmd.Flags().StringVarP(&cvOpt.To, "to", "t", "", "target format (zstd, zip, tar, 7z)")
	cvCmd.Flags().StringVarP(&cvOpt.Prefix, "prefix", "p", "", "output filename prefix")
	cvCmd.Flags().StringVarP(&cvOpt.Suffix, "suffix", "s", "", "output filename suffix")
	cvCmd.Flags().StringVarP(&cvOpt.OutDir, "output-dir", "o", "", "output directory")
	cvCmd.Flags().BoolVarP(&cvOpt.DeleteOriginal, "delete", "d", false, "delete original after conversion")
	cvCmd.Flags().BoolVar(&cvOpt.DryRun, "dry-run", false, "preview output path only")
	rootCmd.AddCommand(cvCmd)
}
