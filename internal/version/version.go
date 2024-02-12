package version

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gitver/internal/gitops"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	PriorityNone = iota
	PriorityFix
	PriorityFeat
	PriorityBreakingChange
)

type IGitVer interface {
	GetVersion() string
	ReadFile() error
	SetConfigFile(path string)
	SafeWriteFile() error
	WriteFile() error
	BumpMajor() error
	BumpMinor() error
	BumpPatch() error
	BumpAuto() error
	BumpFlex() error
	SetFileSystem(fs afero.Fs)
	GetLastVersion() string
	getDebug() bool
	SetDebug(debug bool)
	Commit(amend bool) error
	Push() error
	Add() error
	Tag() error
}

type GitVer struct {
	lastVersion string
	viper       *viper.Viper
	vcs         gitops.IVcs
	debug       bool
}

type Version struct {
	major int
	minor int
	patch int
}

var g IGitVer

func (v *Version) FromString(value string) error {
	_, err := fmt.Sscanf(value, "%d.%d.%d", &v.major, &v.minor, &v.patch)
	if err != nil {
		return err
	}
	return nil
}
func (v *Version) ToString() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}
func init() {
	g = build()
}
func (g *GitVer) GetVersion() string {
	return g.viper.GetString("release.version")
}
func (g *GitVer) SetFileSystem(fs afero.Fs) {
	g.viper.SetFs(fs)
}
func (g *GitVer) ReadFile() error {
	return g.viper.ReadInConfig()
}
func (g *GitVer) SetConfigFile(path string) {
	g.viper.SetConfigFile(path)
}
func (g *GitVer) SafeWriteFile() error {
	return viper.SafeWriteConfig()
}
func (g *GitVer) WriteFile() error {
	return g.viper.WriteConfig()
}
func (g *GitVer) BumpMajor() error {
	verString := g.viper.GetString("release.version")
	version := new(Version)
	err := version.FromString(verString)
	if err != nil {
		return err
	}
	g.lastVersion = version.ToString()
	version.major++
	version.minor = 0
	version.patch = 0
	g.viper.Set("release.version", version.ToString())
	return nil
}
func (g *GitVer) BumpMinor() error {
	verString := g.viper.GetString("release.version")
	version := new(Version)
	err := version.FromString(verString)
	if err != nil {
		return err
	}
	g.lastVersion = version.ToString()
	version.minor++
	version.patch = 0
	g.viper.Set("release.version", version.ToString())
	return nil
}
func (g *GitVer) BumpPatch() error {
	verString := g.viper.GetString("release.version")
	version := new(Version)
	err := version.FromString(verString)
	if err != nil {
		return err
	}
	g.lastVersion = version.ToString()
	version.patch++
	g.viper.Set("release.version", version.ToString())
	return nil
}
func (g *GitVer) BumpAuto() error {

	tag, err := g.vcs.GetLatestTag()
	if err != nil {
		return err
	}

	commits, err := g.vcs.GetCommitsMessagesFromTagToHead(tag)
	if err != nil {
		log.Fatal(err)
	}

	priority := g.findCommitPriority(commits)

	switch priority {
	case PriorityBreakingChange:
		return g.BumpMajor()
	case PriorityFeat:

		return g.BumpMinor()
	case PriorityFix:
		return g.BumpPatch()
	default:
		return nil
	}
}
func (g *GitVer) BumpFlex() error {

	var majorCommand = g.viper.GetString("config.modes.flex.trigger.major")
	var minorCommand = g.viper.GetString("config.modes.flex.trigger.minor")
	var patchCommand = g.viper.GetString("config.modes.flex.trigger.patch")
	var autoCommand = g.viper.GetString("config.modes.flex.trigger.auto")

	commit, err := g.vcs.GetHeadCommit()
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(commit, majorCommand) {
		return g.BumpMajor()
	} else if strings.Contains(commit, minorCommand) {
		return g.BumpMinor()
	} else if strings.Contains(commit, patchCommand) {
		return g.BumpPatch()
	} else if strings.Contains(commit, autoCommand) {
		return g.BumpAuto()
	} else {
		return nil
	}
}
func (g *GitVer) GetLastVersion() string {
	return g.lastVersion
}
func (g *GitVer) getDebug() bool {
	return g.debug
}
func (g *GitVer) SetDebug(debug bool) {
	g.debug = debug
}
func (g *GitVer) Add() error {
	return g.vcs.AddAll()
}
func (g *GitVer) Commit(amend bool) error {
	var message = g.viper.GetString("config.vcs.messages.commit")
	return g.vcs.Commit(message, amend)
}
func (g *GitVer) Push() error {
	return g.vcs.Push()
}
func (g *GitVer) Tag() error {
	var message = g.viper.GetString("config.vcs.messages.tag")
	var format = g.viper.GetString("config.vcs.format.versionTag")
	return g.vcs.CreateTag(fmt.Sprintf(format, g.GetVersion()), message)
}
func (g *GitVer) prepareGitOperation() error {

	if err := g.vcs.ReadRepository(); err != nil {
		return err
	}

	isClean, err := g.vcs.IsCleanRepo()
	if err != nil {
		return err
	}

	if !isClean {
		return fmt.Errorf("")
	}

	return nil
}
func (g *GitVer) findCommitPriority(commits []string) int {
	highestPriority := PriorityNone
	var majorTrigger = g.viper.GetString("config.modes.auto.trigger.major")
	var minorTrigger = g.viper.GetString("config.modes.auto.trigger.minor")
	var patchTrigger = g.viper.GetString("config.modes.auto.trigger.patch")

	for _, commit := range commits {
		message := commit

		if strings.Contains(message, majorTrigger) {
			return PriorityBreakingChange
		} else if strings.Contains(message, minorTrigger) && highestPriority < PriorityFeat {
			highestPriority = PriorityFeat
		} else if strings.Contains(message, patchTrigger) && highestPriority < PriorityFix {
			highestPriority = PriorityFix
		}
	}

	return highestPriority
}
func findProjectDir(fs afero.Fs, dir, folderName string) (string, error) {

	for {

		exist, err := afero.Exists(fs, filepath.Join(dir, folderName))
		if err != nil {
			return "", err
		}
		if exist {
			return filepath.Clean(dir), nil
		}

		// Erhalten Sie das übergeordnete Verzeichnis.
		parentDir := filepath.Dir(dir)

		// Wenn das aktuelle Verzeichnis gleich dem übergeordneten Verzeichnis ist,
		// dann haben wir das Wurzelverzeichnis erreicht.
		if parentDir == dir {
			break
		}

		dir = parentDir
	}

	return "", DirectoryNotFoundError(dir)
}

func Get() IGitVer {
	if g == nil {
		g = build()
	}
	return g
}
func build() IGitVer {
	v := new(GitVer)
	fs := afero.NewOsFs()
	wd, err := os.Getwd()
	if err != nil {

	}
	projectDir, err := findProjectDir(fs, wd, ".gitver")
	if err != nil {

	}

	v.viper = viper.GetViper()
	v.SetFileSystem(fs)
	v.SetConfigFile(filepath.Join(projectDir, ".version.yaml"))
	v.viper.Set("config.default.version", "0.0.0")
	return v
}
func logf(format string, v ...any) {
	if g.getDebug() {
		log.Printf(format, v...)
	}
}
func prints(v ...any) {
	if g.getDebug() {
		log.Print(v...)
	}
}
