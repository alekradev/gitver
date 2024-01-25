package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"gotver/internal/constants"
	"gotver/internal/gitops"
	"gotver/internal/version"
	"log"
	"strings"
)

var (
	autoFlag   bool
	majorFlag  bool
	minorFlag  bool
	patchFlag  bool
	commitFlag bool
	gitFlag    string
	amend      bool

	verbose bool
)

const (
	CommitTag     = "COMMIT_TAG"
	CommitTagPush = "COMMIT_TAG_PUSH"
)

const (
	PriorityNone = iota
	PriorityFix
	PriorityFeat
	PriorityBreakingChange
)

const (
	message0001 = "Please provide a valid flag: --auto, --commit, --major, --minor, or --patch"
	message0002 = "Version bumped: %v -> %v"
)

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Bump the version of the project",
	Long:  "Bump the version of the project based on the provided flags: --auto, --commit, --major, --minor, or --patch.",
	Run: func(cmd *cobra.Command, args []string) {
		if !validateMode(majorFlag, minorFlag, patchFlag, autoFlag, commitFlag) {
			log.Fatalf(message0001)
		}

		loadConfig()

		if gitFlag == CommitTag || gitFlag == CommitTagPush || autoFlag {
			err := prepareGitOperation()
			if err != nil {
				log.Fatal(err)
			}

		}

		switch {
		case majorFlag:
			executeMajorMode()
		case minorFlag:
			executeMinorMode()
		case patchFlag:
			executePatchMode()
		case autoFlag:
			executeAutoMode()
		case commitFlag:
			executeCommitMode()
		}

		executeGitOperations()

		log.Printf(message0002, version.GetLastVersion(), version.ToString())
	},
}

func init() {
	rootCmd.AddCommand(bumpCmd)
	bumpCmd.Flags().BoolVar(&autoFlag, "auto", false, "Bump version from git commits")
	bumpCmd.Flags().BoolVar(&commitFlag, "commit", false, "Bump version from git commits")
	bumpCmd.Flags().BoolVar(&majorFlag, "major", false, "Bump the major version")
	bumpCmd.Flags().BoolVar(&minorFlag, "minor", false, "Bump the minor version")
	bumpCmd.Flags().BoolVar(&patchFlag, "patch", false, "Bump the patch version")
	bumpCmd.Flags().BoolVar(&verbose, "verbose", false, "Bump the patch version")
	bumpCmd.Flags().BoolVar(&amend, "amend", false, "Bump the patch version")
	bumpCmd.Flags().StringVar(&gitFlag, "git", "", "Auto Commit Bump version changes. Valid Values are COMMIT_TAG COMMIT_TAG_PUSH")
}

func validateMode(values ...bool) bool {
	trueCount := 0
	for _, value := range values {
		if value {
			trueCount++
		}
	}
	return trueCount == 1
}

func executeMajorMode() {
	prints("bump major version")
	err := version.BumpMajor()
	if err != nil {
		log.Fatal(err)
	}
	prints("bump major version success")
}

func executeMinorMode() {
	prints("bump minor version")
	err := version.BumpMinor()
	if err != nil {
		log.Fatal(err)
	}
	prints("bump minor version success")
}

func executePatchMode() {
	prints("bump patch version success")
	err := version.BumpPatch()
	if err != nil {
		log.Fatal(err)
	}
	prints("bump patch version success")
}

func executeAutoMode() {
	prints("start auto mode")
	bumpFunc, err := detectAutoBump()
	if err != nil {
		log.Fatal(err)
	}
	err = bumpFunc()
	if err != nil {
		log.Fatal(err)
	}
	prints("start auto mode success")
}

func executeCommitMode() {
	prints("start commit mode")
	bumpFunc, err := detectCommitBump()
	if err != nil {
		log.Fatal(err)
	}
	err = bumpFunc()
	if err != nil {
		log.Fatal(err)
	}
	prints("start commit mode success")
}

func detectAutoBump() (func() error, error) {
	prints("start detect bump function for auto mode")
	tag, err := gitops.GetLastTag()
	if err != nil {
		return nil, err
	}

	if tag == fmt.Sprintf(constants.ReleaseTag, version.ToString()) {
		return analyzeCommits(tag)
	} else if tag == fmt.Sprintf(constants.VersionTag, version.ToString()) {
		return analyzeAndCompareCommits(tag, fmt.Sprintf(constants.ReleaseTag, version.GetLastVersion()))
	} else if tag == "" {
		return analyzeCommits(tag)
	} else {
		return nil, fmt.Errorf("no new version required")
	}
}

func analyzeCommits(tag string) (func() error, error) {
	logf("analyze commits from head to %s...", tag)
	commits, err := gitops.GetCommits(tag)
	if err != nil {
		log.Fatal(err)
	}

	logf("%d commits found", len(commits))

	priority := findCommitPriority(commits)
	logf("bump priority are: %d", priority)
	switch priority {
	case PriorityBreakingChange:
		return version.BumpMajor, nil
	case PriorityFeat:
		return version.BumpMinor, nil
	case PriorityFix:
		return version.BumpPatch, nil
	default:
		return nil, fmt.Errorf("no new version required")
	}
}

func analyzeAndCompareCommits(starttag, endtag string) (func() error, error) {
	logf("analyze commits from %s to %s...", starttag, endtag)
	oldCommits, err := gitops.GetCommitsBetweenTags(starttag, endtag)
	if err != nil {
		log.Fatal(err)
	}
	logf("%d commits found", len(oldCommits))

	logf("analyze commits from head to %s...", starttag)
	newCommits, err := gitops.GetCommits(starttag)
	if err != nil {
		log.Fatal(err)
	}
	logf("%d commits found", len(newCommits))

	priority := comparePriority(findCommitPriority(oldCommits), findCommitPriority(newCommits))
	logf("bump priority are: %d", priority)
	switch priority {
	case PriorityBreakingChange:
		return version.BumpMajor, nil
	case PriorityFeat:
		return version.BumpMinor, nil
	case PriorityFix:
		return version.BumpPatch, nil
	default:
		return nil, fmt.Errorf("cannot detect bump function")
	}
}

func comparePriority(first, second int) int {
	if second > first {
		return second
	}
	return 0
}

func detectCommitBump() (func() error, error) {
	commit, err := gitops.GetHeadCommit()
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(commit.Message, "[bump]") {
		return detectAutoBump()
	} else if strings.Contains(commit.Message, "[bump major]") {
		return version.BumpMajor, nil
	} else if strings.Contains(commit.Message, "[bump minor]") {
		return version.BumpMinor, nil
	} else if strings.Contains(commit.Message, "[bump patch]") {
		return version.BumpPatch, nil
	}

	return nil, fmt.Errorf("no bump commanded")
}

func executeGitOperations() {
	if gitFlag == CommitTag || gitFlag == CommitTagPush {

		if _, err := gitops.Add(); err != nil {
			log.Fatal(err)
		}

		if err := gitops.Commit(fmt.Sprintf(constants.CommitMessage, version.GetLastVersion(), version.ToString()), amend); err != nil {
			log.Fatal(err)
		}

		if err := gitops.CreateTag(fmt.Sprintf(constants.VersionTag, version.ToString()), constants.TagMessage); err != nil {
			log.Fatal(err)
		}
	}

	if gitFlag == CommitTagPush {
		if err := gitops.Push(); err != nil {
			log.Fatal(err)
		}
	}
}

func findCommitPriority(commits []*object.Commit) int {
	highestPriority := PriorityNone

	for _, commit := range commits {
		message := commit.Message

		if strings.Contains(message, "BREAKING CHANGE:") {
			logf("BREAKING CHANGE found: %s", commit.Hash)
			return PriorityBreakingChange
		} else if strings.Contains(message, "feat:") && highestPriority < PriorityFeat {
			logf("FEATURE found: %s", commit.Hash)
			highestPriority = PriorityFeat
		} else if strings.Contains(message, "fix:") && highestPriority < PriorityFix {
			logf("FIX found: %s", commit.Hash)
			highestPriority = PriorityFix
		}
	}

	return highestPriority
}
