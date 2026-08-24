package cli

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"komp/internal/ui"
)

func runInteractive(cmd *cobra.Command) error {
	for {
		action, err := ui.PickCommand()
		if err != nil {
			if ui.IsAbort(err) {
				return nil
			}
			return cliErr(err, 1)
		}
		err = dispatchInteractive(cmd, action)
		if err == nil {
			return nil
		}
		if ui.IsAbort(err) || errors.Is(err, ui.ErrBack) {
			continue
		}
		return err
	}
}

func dispatchInteractive(cmd *cobra.Command, action string) error {
	switch action {
	case "compress":
		return runCreate(cmd, nil)
	case "add":
		archive, err := ui.PickArchiveFor(ui.OpAdd)
		if err != nil {
			return err
		}
		return addCmd.RunE(addCmd, []string{archive})
	case "ls":
		archive, err := ui.PickArchiveFor(ui.OpList)
		if err != nil {
			return err
		}
		return lsCmd.RunE(lsCmd, []string{archive})
	case "un":
		archive, err := ui.PickArchiveFor(ui.OpExtract)
		if err != nil {
			return err
		}
		return unCmd.RunE(unCmd, []string{archive})
	case "clean":
		archive, groups, err := ui.PickCleanSettings()
		if err != nil {
			return err
		}
		if len(groups) == 0 {
			return cliErr(errors.New("no junk groups selected"), 1)
		}
		if err := cleanCmd.Flags().Set("groups", strings.Join(groups, ",")); err != nil {
			return cliErr(err, 1)
		}
		return cleanCmd.RunE(cleanCmd, []string{archive})
	case "t":
		archive, err := ui.PickArchiveFor(ui.OpTest)
		if err != nil {
			return err
		}
		return tCmd.RunE(tCmd, []string{archive})
	case "cv":
		archive, to, err := ui.PickConvertSettings()
		if err != nil {
			return err
		}
		if err := cvCmd.Flags().Set("to", to); err != nil {
			return cliErr(err, 1)
		}
		return cvCmd.RunE(cvCmd, []string{archive})
	case "img":
		src, kind, volname, err := ui.PickImageSettings()
		if err != nil {
			return err
		}
		if err := imgCmd.Flags().Set("type", kind); err != nil {
			return cliErr(err, 1)
		}
		if volname != "" {
			if err := imgCmd.Flags().Set("volname", volname); err != nil {
				return cliErr(err, 1)
			}
		}
		return imgCmd.RunE(imgCmd, []string{src})
	}
	return nil
}
