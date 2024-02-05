package gitops

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gitver/internal/constants"
	"sort"
	"time"
)

const (
	EMPTY = ""
)

var g *GitOps

type GitOps struct {
	path          string
	name          string
	email         string
	repository    *git.Repository
	worktree      *git.Worktree
	object        Object
	commitMessage string
	tagMessage    string
}

// Private Functions

func init() {
	g = New()
}

func (g *GitOps) setRepositoryPath(path string) {
	g.path = path
}

func (g *GitOps) setObject(object Object) {
	g.object = object
}

func (g *GitOps) setWorktree(w *git.Worktree) {
	g.worktree = w
}

func (g *GitOps) setRepository(r *git.Repository) {
	g.repository = r
}

func (g *GitOps) readRepository() error {

	r, err := g.object.Open(g.path)
	if err != nil {
		// RepositoryNotExists
		return err
	}

	cfg, err := r.ConfigScoped(config.GlobalScope)
	if err != nil {
		return err
	}

	user := cfg.User
	if user.Name == "" && user.Email == "" {
		return nil
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	g.name = user.Name
	g.email = user.Email
	g.repository = r
	g.worktree = w

	return nil
}

func (g *GitOps) getStatus() (git.Status, error) {
	return g.worktree.Status()
}
func (g *GitOps) isCleanRepo() (bool, error) {
	status, err := g.getStatus()
	if err != nil {
		return false, err
	}
	return status.IsClean(), nil
}
func (g *GitOps) addAll() (plumbing.Hash, error) {
	return g.worktree.Add(".")
}
func (g *GitOps) commit(message string, amend bool) error {
	commitOptions := &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.name,
			Email: g.email,
			When:  time.Now(),
		},
		Amend: amend,
	}

	if amend {
		commit, err := g.getHeadCommit()
		if err != nil {
			return err
		}
		message = commit.Message
	}

	_, err := g.worktree.Commit(message, commitOptions)

	if err != nil {
		return err
	}
	return nil
}
func (g *GitOps) createTag(tag string, message string) error {

	headRef, err := g.repository.Head()
	if err != nil {
		return err
	}

	headCommit, err := g.repository.CommitObject(headRef.Hash())
	if err != nil {
		return err
	}

	// Create the tag
	_, err = g.repository.CreateTag(tag, headCommit.Hash, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  g.name,
			Email: g.email,
			When:  time.Now(),
		},
		Message: message,
	})
	if err != nil {
		return err
	}

	return nil
}
func (g *GitOps) getLatestTag() (string, error) {

	tagRefs, err := g.repository.Tags()
	if err != nil {
		return "", err
	}

	var tags []struct {
		Name string
		When time.Time
	}

	err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		obj, err := g.repository.TagObject(t.Hash())
		if err != nil {
			// Es könnte ein leichtes Tag sein, also versuchen Sie, den commit direkt zu bekommen
			commit, err := g.repository.CommitObject(t.Hash())
			if err != nil {
				return nil // Ignorieren von Fehlern, die durch leichte Tags verursacht werden
			}
			tags = append(tags, struct {
				Name string
				When time.Time
			}{t.Name().Short(), commit.Committer.When})
			return nil
		}
		tags = append(tags, struct {
			Name string
			When time.Time
		}{t.Name().Short(), obj.Tagger.When})
		return nil
	})
	if err != nil {
		return "", err
	}

	// Sortieren Sie die Tags nach Datum
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].When.After(tags[j].When)
	})

	if len(tags) > 0 {
		return tags[0].Name, nil
	}

	return "", nil
}
func (g *GitOps) push() error {

	// HEAD-Referenz holen
	headRef, err := g.repository.Head()
	if err != nil {
		return err
	}

	// HEAD dereferenzieren, um den aktuellen Branch zu bekommen
	ref, err := g.repository.Reference(headRef.Name(), true)
	if err != nil {
		return err
	}

	err = g.repository.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec(ref.Name() + ":" + ref.Name())},
	})
	if err != nil {
		return err
	}

	return nil
}
func (g *GitOps) getHeadCommit() (*object.Commit, error) {
	headRef, err := g.repository.Head()
	if err != nil {
		return nil, err
	}

	commit, err := g.repository.CommitObject(headRef.Hash())
	if err != nil {

	}

	return commit, nil
}
func (g *GitOps) getTag(tag string) (plumbing.Hash, error) {
	tagRef, err := g.repository.Tag(tag)
	if err == nil {
		// Dereferenziert das Tag-Objekt, falls es ein annotiertes Tag ist
		resolvedTag, err := g.repository.TagObject(tagRef.Hash())
		if err == nil {
			return resolvedTag.Target, nil
		} else {
			// Wenn es kein annotiertes Tag ist, sondern ein leichtgewichtiger Tag
			return tagRef.Hash(), nil
		}
	}
	return plumbing.Hash{}, nil
}
func (g *GitOps) getCommitsBetweenTags(startTag, endTag string) ([]*object.Commit, error) {

	startHash, err := g.getTag(startTag)
	if err != nil {
		headRef, err := g.repository.Head()
		if err != nil {
			return nil, err
		}
		startHash = headRef.Hash()
	}

	endHash, _ := g.getTag(endTag)
	if err != nil {
		endHash = plumbing.Hash{}
	}

	logOptions := &git.LogOptions{From: startHash}
	commitIter, err := g.repository.Log(logOptions)
	if err != nil {
		return nil, err
	}

	var commits []*object.Commit
	err = commitIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		if c.Hash == endHash {
			return fmt.Errorf("BREAK")
		}
		return nil
	})

	if err != nil && err.Error() != "BREAK" {
		return nil, err
	}

	return commits, nil
}

// Public Functions

func New() *GitOps {
	g := new(GitOps)
	g.commitMessage = constants.CommitMessage
	g.tagMessage = constants.TagMessage
	return g
}
func SetRepositoryPath(path string) {
	g.setRepositoryPath(path)
}
func ReadRepository() error {
	return g.readRepository()
}
func GetStatus() (git.Status, error) {
	return g.getStatus()
}
func IsCleanRepo() (bool, error) {
	return g.isCleanRepo()
}
func AddAll() (plumbing.Hash, error) {
	return g.addAll()
}
func Commit(message string, amend bool) error {
	return g.commit(message, amend)
}
func CreateTag(tag string, message string) error {
	return g.createTag(tag, message)
}
func GetLastTag() (string, error) {
	return g.getLatestTag()
}
func Push() error {
	return g.push()
}
func GetHeadCommit() (*object.Commit, error) {
	return g.getHeadCommit()
}
func GetCommits(tag string) ([]*object.Commit, error) {
	return g.getCommitsBetweenTags(EMPTY, tag)
}
func GetTag(tag string) (plumbing.Hash, error) {
	return g.getTag(tag)
}
func HasTag(tag string) bool {
	hash, _ := g.getTag(tag)
	return !hash.IsZero()
}
func GetCommitsBetweenTags(tagStart, tagEnd string) ([]*object.Commit, error) {
	return g.getCommitsBetweenTags(tagStart, tagEnd)
}
