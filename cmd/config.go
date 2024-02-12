package cmd

import (
	"github.com/spf13/cobra"
)

const (
	DescConfigCommand      = "config"
	DescConfigCommandShort = ""
	DescConfigCommandLong  = ""
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   DescConfigCommand,
	Short: DescConfigCommandShort,
	Long:  DescConfigCommandLong,
}

func init() {
	//rootCmd.AddCommand(configCmd)
}
