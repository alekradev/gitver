package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

const (
	DescMajorCommand      = "major"
	DescMajorCommandShort = "Executes a major version bump."
	DescMajorCommandLong  = "The major command is used to manually increment the major version of your project, signifying incompatible API changes or substantial modifications. This command should be used when introducing changes that are not backward compatible, ensuring that the versioning clearly communicates the potential impact on existing users or dependent systems.\n\n"
	MsgBumpMajorInit      = "Initiating major version bump..."
	MsgBumpMajorSuccess   = "Major version bumped successfully."
)

// majorCmd represents the major command
var majorCmd = &cobra.Command{
	Use:   DescMajorCommand,
	Short: DescMajorCommandShort,
	Long:  DescMajorCommandLong,
	Run:   runMajorCmd,
}

func init() {
	releaseCmd.AddCommand(majorCmd)
}

func runMajorCmd(cmd *cobra.Command, args []string) {
	prints(MsgBumpMajorInit)
	err := F.BumpMajor()
	if err != nil {
		log.Fatal(err)
	}
	prints(MsgBumpMajorSuccess)
}
