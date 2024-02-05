/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitver/internal/version"
	"log"
)

var (
	versionFlag string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: executeInitCmd,
}

func init() {
	configCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&versionFlag, "version", "v", "0.0.1", "overrides the default initial version")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func executeInitCmd(cmd *cobra.Command, args []string) {
	err := version.SetDefaultVersion(versionFlag)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = version.SafeWriteFile()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	err = viper.SafeWriteConfig()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	log.Println("Gitver initialized for the project.")
}
