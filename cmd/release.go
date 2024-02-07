package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitver/internal/constants"
	"log"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   DescReleaseCommand,
	Short: DescReleaseCommandShort,
	Long:  DescReleaseCommandLong,
	Run:   executeReleaseCmd,
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}

func executeReleaseCmd(cmd *cobra.Command, args []string) {
	loadConfig()
	err := prepareGitOperation()
	if err != nil {
		log.Fatal(err)
	}

	if err := G.CreateTag(fmt.Sprintf(constants.ReleaseTag, F.GetVersion()), constants.TagMessage); err != nil {
		return
	}
	log.Printf("Tag: %s tagged", fmt.Sprintf(constants.ReleaseTag, F.GetVersion()))
}
