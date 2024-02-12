package cmd

import (
	"github.com/spf13/cobra"
	"gitver/internal/version"
	"log"
)

const (
	DescReleaseCommand      = "bump"
	DescReleaseCommandShort = "Bump the project version."
	DescReleaseCommandLong  = "Bumps the project version based on the specified flags: --auto, --commit, --major, --minor, or --patch."
	DescAddFlag             = "Add all modified files to the stage"
	DescCommitFlag          = "Creates a commit"
	DescAmendFlag           = "Amends the last commit during the version bump process."
	DescTagFlag             = "Creates a tag for the new version."
	DescPushFlag            = "Pushes changes to the remote git repository."
	DescPrereleaseFlag      = "Create a version based on the modus but will not bumped. Excluded git operations"
	NameAddFlag             = "add"
	NameAmendFlag           = "amend"
	NameCommitFlag          = "commit"
	NameTagFlag             = "tag"
	NamePushFlag            = "push"
	NamePrereleaseFlag      = "prerelease"

	ErrWriteFile                  = "error by write file. causes: %s"
	ErrGitAddFailed               = "error by add files to git stage. causes: %s"
	ErrGitCommitFailed            = "error by commit. causes: %s"
	ErrGitTagFailed               = "error by create git tag. causes: %s"
	ErrGitPushFailed              = "error by push git repository. causes: %s"
	ErrPrereleaseWithGitOperation = "the flag --prerelease is not allowed with --add, --commit, --tag, --push"
	ErrAmendWithoutTag            = "--amend flag requires --commit flag to be set"
)

var (
	commitFlag     bool
	amendFlag      bool
	verboseFlag    bool
	pushFlag       bool
	tagFlag        bool
	addFlag        bool
	prereleaseFlag bool
)

// releaseCmd represents the bump command
var releaseCmd = &cobra.Command{
	Use:               DescReleaseCommand,
	Short:             DescReleaseCommandShort,
	Long:              DescReleaseCommandLong,
	PersistentPreRun:  preRunReleaseCmd,
	PersistentPostRun: postRunReleaseCmd,
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.PersistentFlags().BoolVar(&addFlag, NameAddFlag, false, DescAddFlag)
	releaseCmd.PersistentFlags().BoolVar(&amendFlag, NameAmendFlag, false, DescAmendFlag)
	releaseCmd.PersistentFlags().BoolVar(&commitFlag, NameCommitFlag, false, DescCommitFlag)
	releaseCmd.PersistentFlags().BoolVar(&tagFlag, NameTagFlag, false, DescTagFlag)
	releaseCmd.PersistentFlags().BoolVar(&pushFlag, NamePushFlag, false, DescPushFlag)
	releaseCmd.PersistentFlags().BoolVar(&prereleaseFlag, NamePrereleaseFlag, false, DescPrereleaseFlag)
}

func preRunReleaseCmd(cmd *cobra.Command, args []string) {
	if prereleaseFlag && (amendFlag || commitFlag || tagFlag || pushFlag || addFlag) {
		log.Fatal(ErrPrereleaseWithGitOperation)
	}

	if !commitFlag && amendFlag {
		log.Fatal(ErrAmendWithoutTag)
	}
}

func postRunReleaseCmd(cmd *cobra.Command, args []string) {
	var v = version.Get()
	if !prereleaseFlag {
		err := v.SafeWriteFile()
		if err != nil {
			log.Fatalf(ErrWriteFile, err)
		}

		if addFlag {
			err := v.Add()
			if err != nil {
				log.Printf(ErrGitAddFailed, err)
			}
		}
		if commitFlag {
			err := v.Commit(amendFlag)
			if err != nil {
				log.Printf(ErrGitCommitFailed, err)
			}
		}

		if tagFlag {
			err := v.Tag()
			if err != nil {
				log.Printf(ErrGitTagFailed, err)
			}
		}

		if pushFlag {
			err := v.Push()
			if err != nil {
				log.Printf(ErrGitPushFailed, err)
			}
		}
	}
}
