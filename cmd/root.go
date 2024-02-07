package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"gitver/internal/constants"
	"gitver/internal/filesystem"
	"gitver/internal/gitops"
	"gitver/internal/version"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const (
	MsgBumpMajorInit       = "Initiating major version bump..."
	MsgBumpMajorSuccess    = "Major version bumped successfully."
	MsgBumpMinorInit       = "Initiating minor version bump..."
	MsgBumpMinorSuccess    = "Minor version bumped successfully."
	MsgBumpPatchInit       = "Initiating patch version bump..."
	MsgBumpPatchSuccess    = "Patch version bumped successfully."
	MsgAutoBumpInit        = "Initiating automatic version bump..."
	MsgAutoBumpSuccess     = "Automatic version bump completed successfully."
	MsgCommitBumpInit      = "Detecting version bump from commit message..."
	MsgCommitBumpSuccess   = "Version bump from commit message completed successfully."
	MsgVersionChange       = "Version changed from %v to %v."
	MsgDetectingBumpMethod = "Detecting bump method for automatic mode..."

	MsgBreakingChangeDetected = "Detected BREAKING CHANGE in commit: %s"
	MsgFeatureDetected        = "Detected new feature in commit: %s"
	MsgFixDetected            = "Detected fix in commit: %s"
	MsgBumpPrioritySet        = "Bump priority set to: %d"
	MsgPreparingGitOperation  = "Preparing Git operation..."
	MsgConfigLoaded           = "Configuration loaded."
	MsgGitOperationsCompleted = "Git operations completed successfully."

	ErrInvalidModeFlag             = "invalid mode flag provided. Please use one of --auto, --commit, --major, --minor, or --patch"
	ErrPushWithoutTag              = "--push flag requires --tag flag to be set"
	ErrAmendWithoutTag             = "--amend flag requires --tag flag to be set"
	ErrNoVersionBump               = "no version bump required based on the current criteria"
	ErrGitOperationFailed          = "git operation failed"
	ErrBumpMethodDetection         = "failed to detect bump method"
	ErrNoBumpCommandFound          = "no bump command found in commit message"
	ErrGitRepositoryNotInitialized = "git repository is not initialized %w"
	ErrWorktreeIsNotClean          = "working tree is not clean"

	LogAnalyzingCommits = "Analyzing commits from %s to %s..."
	LogCommitsFound     = "%d commits found."

	DescAutoFlag    = "Automatically determines and sets the new version based on git commits."
	DescCommitFlag  = "Sets the new version based on keywords in commit messages."
	DescMajorFlag   = "Increments the major version."
	DescMinorFlag   = "Increments the minor version."
	DescPatchFlag   = "Increments the patch version."
	DescVerboseFlag = "Enables verbose output."
	DescAmendFlag   = "Amends the last commit during the version bump process."
	DescTagFlag     = "Creates a commit and tag for the new version."
	DescPushFlag    = "Pushes changes to the remote git repository."

	DescBumpCommand      = "bump"
	DescBumpCommandShort = "Bump the project version."
	DescBumpCommandLong  = "Bumps the project version based on the specified flags: --auto, --commit, --major, --minor, or --patch."

	DescRootCommand      = "gitver"
	DescRootCommandShort = ""
	DescRootCommandLong  = ""

	DescInitCommand      = "init"
	DescInitCommandShort = ""
	DescInitCommandLong  = ""

	DescConfigCommand      = "config"
	DescConfigCommandShort = ""
	DescConfigCommandLong  = ""

	DescReleaseCommand      = "release"
	DescReleaseCommandShort = ""
	DescReleaseCommandLong  = ""
)

type cmd struct {
	v *viper.Viper
	f version.IFile
	g gitops.IGitOps
}

var V *viper.Viper
var F version.IFile
var G gitops.IGitOps

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   DescRootCommand,
	Short: DescRootCommandShort,
	Long:  DescRootCommandLong,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func (c *cmd) init() {
	V = viper.GetViper()
	F = version.Get()
	G = gitops.Get()

	projectDir, err := filesystem.New().Find(constants.ConfigFolderName)
	if err != nil {
		projectDir, err = os.Getwd()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	V.SetConfigName(constants.ConfigName)
	V.SetConfigType(constants.ConfigType)
	V.AddConfigPath(projectDir + "/" + constants.ConfigFolderName)
	V.SetDefault("data", "0.0.0")

	F.SetFilePath(projectDir + "/" + constants.ConfigFolderName)
	F.SetFileName(constants.VersionFileName)

	G.SetRepositoryPath(projectDir)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func init() {

}

func loadConfig() {

	if err := V.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := F.ReadFile(); err != nil {
		log.Fatal(err)
	}
	prints(MsgConfigLoaded)
}

func prepareGitOperation() error {
	prints(MsgPreparingGitOperation)
	if err := G.ReadRepository(); err != nil {
		return fmt.Errorf(ErrGitRepositoryNotInitialized, err)
	}

	isClean, err := G.IsCleanRepo()
	if err != nil {
		return err
	}

	if !isClean {
		return fmt.Errorf(ErrWorktreeIsNotClean)
	}

	prints(MsgGitOperationsCompleted)
	return nil

}

func logf(format string, v ...any) {
	if verboseFlag {
		log.Printf(format, v...)
	}
}

func prints(v ...any) {
	if verboseFlag {
		log.Print(v...)
	}
}
