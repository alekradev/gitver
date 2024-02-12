package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

const (
	DescFlexCommand      = "flex"
	DescFlexCommandShort = "Dynamically determines version bump from commit messages."
	DescFlexCommandLong  = "The flex command offers a flexible approach to version bumping by analyzing commit messages for specific keywords or patterns that indicate the nature of the changes (e.g., \"feature\" for minor, \"fix\" for patch). This allows for a more nuanced version management strategy, accommodating projects that may not strictly adhere to conventional commit standards or that require custom versioning logic."
	MsgCommitBumpInit    = "Detecting version bump from commit message..."
	MsgCommitBumpSuccess = "Version bump from commit message completed successfully."
)

// flexCmd represents the flex command
var flexCmd = &cobra.Command{
	Use:   DescFlexCommand,
	Short: DescFlexCommandShort,
	Long:  DescFlexCommandLong,
	Run:   runFlexCmd,
}

func init() {
	releaseCmd.AddCommand(flexCmd)
}

func runFlexCmd(cmd *cobra.Command, args []string) {
	prints(MsgCommitBumpInit)
	err := F.BumpFlex()
	if err != nil {
		log.Fatal(err)
	}
	prints(MsgCommitBumpSuccess)
}
