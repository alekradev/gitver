package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

const (
	DescAutoCommand      = "auto"
	DescAutoCommandShort = "Automates version bumping based on pre-defined rules and commit history analysis."
	DescAutoCommandLong  = "The auto command streamlines the version bumping process by automatically determining the appropriate version increment (major, minor, or patch) based on a set of pre-defined rules or by analyzing the commit history. This command is ideal for continuous integration (CI) environments or for developers looking to automate their versioning workflow, ensuring consistent and accurate version updates with minimal manual intervention"

	MsgAutoBumpInit    = "Initiating automatic version bump..."
	MsgAutoBumpSuccess = "Automatic version bump completed successfully."
)

// autoCmd represents the auto command
var autoCmd = &cobra.Command{
	Use:   DescAutoCommand,
	Short: DescAutoCommandShort,
	Long:  DescAutoCommandLong,
	Run:   runAutoCmd,
}

func init() {
	releaseCmd.AddCommand(autoCmd)
}

func runAutoCmd(cmd *cobra.Command, args []string) {
	prints(MsgAutoBumpInit)
	err := F.BumpAuto()
	if err != nil {
		log.Fatal(err)
	}
	prints(MsgAutoBumpSuccess)
}
