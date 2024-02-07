package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	"gitver/internal/constants"
	"gitver/internal/gitops"
	"log"
	"strings"
)

var (
	autoFlag    bool
	majorFlag   bool
	minorFlag   bool
	patchFlag   bool
	commitFlag  bool
	amendFlag   bool
	verboseFlag bool
	pushFlag    bool
	tagFlag     bool
)

const (
	PriorityNone = iota
	PriorityFix
	PriorityFeat
	PriorityBreakingChange
)

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   DescBumpCommand,
	Short: DescBumpCommandShort,
	Long:  DescBumpCommandLong,
	Run:   executeBumpCmd,
}

func init() {
	rootCmd.AddCommand(bumpCmd)
	bumpCmd.Flags().BoolVar(&autoFlag, "auto", false, DescAutoFlag)
	bumpCmd.Flags().BoolVar(&commitFlag, "commit", false, DescCommitFlag)
	bumpCmd.Flags().BoolVar(&majorFlag, "major", false, DescMajorFlag)
	bumpCmd.Flags().BoolVar(&minorFlag, "minor", false, DescMinorFlag)
	bumpCmd.Flags().BoolVar(&patchFlag, "patch", false, DescPatchFlag)
	bumpCmd.Flags().BoolVar(&verboseFlag, "verbose", false, DescVerboseFlag)
	bumpCmd.Flags().BoolVar(&amendFlag, "amend", false, DescAmendFlag)
	bumpCmd.Flags().BoolVar(&tagFlag, "tag", false, DescTagFlag)
	bumpCmd.Flags().BoolVar(&pushFlag, "push", false, DescPushFlag)
}

func validateFlags() (bool, error) {
	if !tagFlag && pushFlag {
		return false, fmt.Errorf(ErrPushWithoutTag)
	}

	if !tagFlag && amendFlag {
		return false, fmt.Errorf(ErrAmendWithoutTag)
	}

	if !validateMode(majorFlag, minorFlag, patchFlag, autoFlag, commitFlag) {
		return false, fmt.Errorf(ErrInvalidModeFlag)
	}

	return true, nil
}

func executeBumpCmd(cmd *cobra.Command, args []string) {

	_, err := validateFlags()
	if err != nil {
		log.Fatal(err)
	}

	loadConfig()

	if tagFlag || pushFlag || autoFlag {
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

	log.Printf(MsgVersionChange, F.GetLastVersion(), F.GetVersion())
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
	prints(MsgBumpMajorInit)
	err := F.BumpMajor()
	if err != nil {
		log.Fatal(err)
	}
	prints(MsgBumpMajorSuccess)
}

func executeMinorMode() {
	prints(MsgBumpMinorInit)
	err := F.BumpMinor()
	if err != nil {
		log.Fatal(err)
	}
	prints(MsgBumpMinorSuccess)
}

func executePatchMode() {
	prints(MsgBumpPatchInit)
	err := F.BumpPatch()
	if err != nil {
		log.Fatal(err)
	}
	prints(MsgBumpPatchSuccess)
}

func executeAutoMode() {
	prints(MsgAutoBumpInit)
	bumpFunc, err := detectAutoBump()
	if err != nil {
		log.Fatal(err)
	}
	err = bumpFunc()
	if err != nil {
		log.Fatal(err)
	}
	prints(MsgAutoBumpSuccess)
}

func executeCommitMode() {
	prints(MsgCommitBumpInit)
	bumpFunc, err := detectCommitBump()
	if err != nil {
		log.Fatal(err)
	}
	err = bumpFunc()
	if err != nil {
		log.Fatal(err)
	}
	prints(MsgCommitBumpSuccess)
}

func detectAutoBump() (func() error, error) {
	prints(MsgDetectingBumpMethod)
	tag, err := G.GetLatestTag()
	if err != nil {
		return nil, err
	}

	if tag == fmt.Sprintf(constants.ReleaseTag, F.GetVersion()) {
		return analyzeCommits(tag)
	} else if tag == fmt.Sprintf(constants.VersionTag, F.GetVersion()) {
		return analyzeAndCompareCommits(tag, fmt.Sprintf(constants.ReleaseTag, F.GetLastVersion()))
	} else if tag == "" {
		return analyzeCommits(tag)
	} else {
		return nil, fmt.Errorf(ErrNoVersionBump)
	}
}

func analyzeCommits(tag string) (func() error, error) {
	logf(LogAnalyzingCommits, "head", tag)
	commits, err := G.GetCommitsBetweenTags(gitops.HEAD, tag)
	if err != nil {
		log.Fatal(err)
	}

	logf(LogCommitsFound, len(commits))

	priority := findCommitPriority(commits)
	logf(MsgBumpPrioritySet, priority)
	switch priority {
	case PriorityBreakingChange:
		return F.BumpMajor, nil
	case PriorityFeat:
		return F.BumpMinor, nil
	case PriorityFix:
		return F.BumpPatch, nil
	default:
		return nil, fmt.Errorf(ErrBumpMethodDetection)
	}
}

func analyzeAndCompareCommits(starttag, endtag string) (func() error, error) {
	logf(LogAnalyzingCommits, starttag, endtag)
	oldCommits, err := G.GetCommitsBetweenTags(starttag, endtag)
	if err != nil {
		log.Fatal(err)
	}
	logf(LogCommitsFound, len(oldCommits))

	logf(LogAnalyzingCommits, "head", starttag)
	newCommits, err := G.GetCommitsBetweenTags(gitops.HEAD, starttag)
	if err != nil {
		log.Fatal(err)
	}
	logf(LogCommitsFound, len(newCommits))

	priority := comparePriority(findCommitPriority(oldCommits), findCommitPriority(newCommits))
	logf(MsgBumpPrioritySet, priority)
	switch priority {
	case PriorityBreakingChange:
		return F.BumpMajor, nil
	case PriorityFeat:
		return F.BumpMinor, nil
	case PriorityFix:
		return F.BumpPatch, nil
	default:
		return nil, fmt.Errorf(ErrNoVersionBump)
	}
}

func comparePriority(first, second int) int {
	if second > first {
		return second
	}
	return 0
}

func detectCommitBump() (func() error, error) {
	commit, err := G.GetHeadCommit()
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(commit.Message, "[bump]") {
		return detectAutoBump()
	} else if strings.Contains(commit.Message, "[bump major]") {
		return F.BumpMajor, nil
	} else if strings.Contains(commit.Message, "[bump minor]") {
		return F.BumpMinor, nil
	} else if strings.Contains(commit.Message, "[bump patch]") {
		return F.BumpPatch, nil
	}

	return nil, fmt.Errorf(ErrNoBumpCommandFound)
}

func executeGitOperations() {
	if tagFlag || pushFlag {

		if _, err := G.AddAll(); err != nil {
			log.Fatal(ErrGitOperationFailed, err)
		}

		if err := G.Commit(fmt.Sprintf(constants.CommitMessage, F.GetLastVersion(), F.GetVersion()), amendFlag); err != nil {
			log.Fatal(ErrGitOperationFailed, err)
		}

		if err := G.CreateTag(fmt.Sprintf(constants.VersionTag, F.GetVersion()), constants.TagMessage); err != nil {
			log.Fatal(ErrGitOperationFailed, err)
		}
	}

	if pushFlag {
		if err := G.Push(); err != nil {
			log.Fatal(ErrGitOperationFailed, err)
		}
	}
}

func findCommitPriority(commits []*object.Commit) int {
	highestPriority := PriorityNone

	for _, commit := range commits {
		message := commit.Message

		if strings.Contains(message, "BREAKING CHANGE:") {
			logf(MsgBreakingChangeDetected, commit.Hash)
			return PriorityBreakingChange
		} else if strings.Contains(message, "feat:") && highestPriority < PriorityFeat {
			logf(MsgFeatureDetected, commit.Hash)
			highestPriority = PriorityFeat
		} else if strings.Contains(message, "fix:") && highestPriority < PriorityFix {
			logf(MsgFixDetected, commit.Hash)
			highestPriority = PriorityFix
		}
	}

	return highestPriority
}
