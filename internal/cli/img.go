package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"komp/internal/image"
)

type imgFlags struct {
	Volname   string
	Format    string
	Size      string
	Fs        string
	ID        string
	Version   string
	SfxModule string
	Out       string
}

var (
	imgOpt  imgFlags
	kindFlag string
)

var imgCmd = &cobra.Command{
	Use:   "img <src>",
	Short: "Build a disk image from a folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		src := args[0]
		kind := kindFlag
		if kind == "" {
			kind = "dmg"
		}
		out := imgOpt.Out
		if out == "" {
			out = src + "." + kind
		}
		if err := image.Build(kind, src, out, image.Vars{
			Volname:  imgOpt.Volname,
			Format:   imgOpt.Format,
			Size:     imgOpt.Size,
			Fs:       imgOpt.Fs,
			ID:       imgOpt.ID,
			Version:  imgOpt.Version,
		}); err != nil {
			return fmt.Errorf("img build: %w", err)
		}
		fmt.Println(out)
		return nil
	},
}

func init() {
	imgCmd.Flags().StringVarP(&imgOpt.Volname, "volname", "", "", "volume name")
	imgCmd.Flags().StringVarP(&imgOpt.Format, "format", "", "", "image format (UDZO, UDZO, etc.)")
	imgCmd.Flags().StringVarP(&imgOpt.Size, "size", "", "", "size for sparse (e.g. 4g)")
	imgCmd.Flags().StringVarP(&imgOpt.Fs, "fs", "", "", "filesystem (APFS, HFS+)")
	imgCmd.Flags().StringVarP(&imgOpt.ID, "id", "", "", "package identifier")
	imgCmd.Flags().StringVarP(&imgOpt.Version, "version", "", "", "package version")
	imgCmd.Flags().StringVarP(&imgOpt.SfxModule, "sfx-module", "", "", "SFX module for self-extracting exe")
	imgCmd.Flags().StringVarP(&imgOpt.Out, "output", "o", "", "output path")
	imgCmd.Flags().StringVarP(&kindFlag, "type", "t", "dmg", "image type: dmg, sparsebundle, sparseimage, iso, pkg")
	rootCmd.AddCommand(imgCmd)
}
