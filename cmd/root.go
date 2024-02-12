package cmd

import (
	"gitver/internal/version"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const (
	MsgVersionChange       = "Version changed from %v to %v."
	MsgDetectingBumpMethod = "Detecting bump method for automatic mode..."

	MsgBreakingChangeDetected = "Detected BREAKING CHANGE in commit: %s"
	MsgFeatureDetected        = "Detected new feature in commit: %s"
	MsgFixDetected            = "Detected fix in commit: %s"
	MsgBumpPrioritySet        = "Bump priority set to: %d"
	MsgPreparingGitOperation  = "Preparing Git operation..."
	MsgConfigLoaded           = "Configuration loaded."
	MsgGitOperationsCompleted = "Git operations completed successfully."

	ErrNoVersionBump               = "no version bump required based on the current criteria"
	ErrGitOperationFailed          = "git operation failed"
	ErrBumpMethodDetection         = "failed to detect bump method"
	ErrNoBumpCommandFound          = "no bump command found in commit message"
	ErrGitRepositoryNotInitialized = "git repository is not initialized %w"
	ErrWorktreeIsNotClean          = "working tree is not clean"

	LogAnalyzingCommits = "Analyzing commits from %s to %s..."
	LogCommitsFound     = "%d commits found."

	DescVerboseFlag = "Enables verbose output."

	DescRootCommand      = "gitver"
	DescRootCommandShort = ""
	DescRootCommandLong  = ""
)

var F version.IGitVer

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

func init() {
	F = version.Get()
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().Bool("verbose", false, DescVerboseFlag)
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
