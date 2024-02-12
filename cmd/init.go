package cmd

import (
	"github.com/spf13/cobra"
	"gitver/internal/version"
	"log"
)

const (
	MsgProjectInitialized = "gitver initialized for the project."
	DescInitCommand       = "init"
	DescInitCommandShort  = "FIND DESCRIPTION"
	DescInitCommandLong   = "FIND DESCRIPTION"

	DescVersionFlag = "overrides the default initial version. Default is 0.0.0"
	NameVersionFlag = "version"

	ErrInitWriteFile = "error by write %s. causes: %s"
)

var (
	versionFlag string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   DescInitCommand,
	Short: DescInitCommandShort,
	Long:  DescInitCommandLong,
	Run:   runInitCmd,
}

func init() {
	configCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&versionFlag, NameVersionFlag, "0.0.0", DescVersionFlag)
}

func runInitCmd(cmd *cobra.Command, args []string) {
	var v = version.Get()

	err := v.SafeWriteFile()
	if err != nil {
		log.Fatalf(ErrInitWriteFile, "", err)
	}

	log.Println(MsgProjectInitialized)
}
