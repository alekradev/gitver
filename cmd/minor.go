package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

const (
	DescMinorCommand      = "minor"
	DescMinorCommandShort = "Performs a minor version bump."
	DescMinorCommandLong  = "The minor command increments the minor version of your project, typically used for adding new functionality in a backward-compatible manner. It's suitable for releases that extend the current functionality without altering the existing API or core behavior, allowing for incremental improvement and expansion of the project's capabilities."

	MsgBumpMinorInit    = "Initiating minor version bump..."
	MsgBumpMinorSuccess = "Minor version bumped successfully."
)

// minorCmd represents the minor command
var minorCmd = &cobra.Command{
	Use:   DescMinorCommand,
	Short: DescMinorCommandShort,
	Long:  DescMinorCommandLong,
	Run:   runMinorCmd,
}

func init() {
	releaseCmd.AddCommand(minorCmd)
}

func runMinorCmd(cmd *cobra.Command, args []string) {
	prints(MsgBumpMinorInit)
	err := F.BumpMinor()
	if err != nil {
		log.Fatal(err)
	}
	prints(MsgBumpMinorSuccess)
}
