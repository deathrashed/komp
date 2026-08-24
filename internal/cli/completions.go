package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "completions [bash|zsh|fish|powershell]",
		Short: "Generate shell completions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				return rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				return rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				return rootCmd.GenPowerShellCompletion(os.Stdout)
			default:
				return fmt.Errorf("unsupported shell %q", args[0])
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "man",
		Short: "Generate a man page",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			man := `.\" komp man page
.TH KOMP 1 "2026-08-25" "komp" "User Commands"
.SH NAME
komp \- unified archive & image toolkit
.SH SYNOPSIS
.B komp
\fIcommand\fR [\fIoptions\fR] [\fIargs...\fR]
.SH COMMANDS
.TP
.B create
Create a new archive.
.TP
.B add
Add files into an existing archive.
.TP
.B ls
List archive contents.
.TP
.B un
Extract an archive.
.TP
.B clean
Strip junk patterns from an archive.
.TP
.B cv
Recompress an archive to another format.
.TP
.B img
Build a disk image from a folder.
.SH OPTIONS
.TP
.B \-\-finder
Use current Finder selection.
.TP
.B \-\-delete
Delete originals after operation.
.TP
.B \-\-dry-run
Preview without changes.
.SH SEE ALSO
.B tar(1), zip(1), 7z(1)
`
			fmt.Fprint(os.Stdout, man)
			return nil
		},
	})
}
