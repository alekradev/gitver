package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	versionFlag string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   DescInitCommand,
	Short: DescInitCommandShort,
	Long:  DescInitCommandLong,
	Run:   executeInitCmd,
}

func init() {
	configCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&versionFlag, "version", "V", "0.0.1", "overrides the default initial version")
}

func executeInitCmd(cmd *cobra.Command, args []string) {
	err := F.SetDefaultVersion(versionFlag)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = F.SafeWriteFile()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	err = V.SafeWriteConfig()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	log.Println("Gitver initialized for the project.")
}
