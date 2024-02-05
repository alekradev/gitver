package cmd

import (
	"github.com/spf13/viper"
	"gitver/internal/constants"
	"gitver/internal/filesystem"
	"gitver/internal/gitops"
	"gitver/internal/version"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   constants.ProgrammName,
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
		examples and usage of using your application. For example:
		
		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gotver.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	projectDir, err := filesystem.New().Find(constants.ConfigFolderName)
	if err != nil {
		projectDir, err = os.Getwd()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	viper.SetConfigName(constants.ConfigName)
	viper.SetConfigType(constants.ConfigType)
	viper.AddConfigPath(projectDir + "/" + constants.ConfigFolderName)
	viper.SetDefault("data", "0.0.0")

	version.SetFilePath(projectDir + "/" + constants.ConfigFolderName)
	version.SetFileName(constants.VersionFileName)

	gitops.SetRepositoryPath(projectDir)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
