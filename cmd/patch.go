package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

const (
	DescPatchCommand      = "patch"
	DescPatchCommandShort = "Initiates a patch version bump."
	DescPatchCommandLong  = "The patch command is designed for making backward-compatible bug fixes, incrementing the patch version of your project. This command is ideal for small updates that rectify errors or issues without introducing new features or making significant changes, ensuring stability and reliability of the software without affecting existing functionality."

	MsgBumpPatchInit    = "Initiating patch version bump..."
	MsgBumpPatchSuccess = "Patch version bumped successfully."
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:   DescPatchCommand,
	Short: DescPatchCommandShort,
	Long:  DescPatchCommandLong,
	Run:   runPatchCmd,
}

func init() {
	releaseCmd.AddCommand(patchCmd)
}

func runPatchCmd(cmd *cobra.Command, args []string) {
	prints(MsgBumpPatchInit)
	err := F.BumpPatch()
	if err != nil {
		log.Fatal(err)
	}
	prints(MsgBumpPatchSuccess)
}
