/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitver/internal/constants"
	"gitver/internal/gitops"
	"gitver/internal/version"
	"log"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: executeReleaseCmd,
}

func init() {
	rootCmd.AddCommand(releaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// releaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// releaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func executeReleaseCmd(cmd *cobra.Command, args []string) {
	loadConfig()
	err := prepareGitOperation()
	if err != nil {
		log.Fatal(err)
	}

	if err := gitops.CreateTag(fmt.Sprintf(constants.ReleaseTag, version.GetVersion()), constants.TagMessage); err != nil {
		return
	}
	log.Printf("Tag: %s tagged", fmt.Sprintf(constants.ReleaseTag, version.GetVersion()))
}
