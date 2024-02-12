package gitops

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"sort"
	"time"
)

type IVcs interface {
	SetRepositoryOpener(object IRepositoryOpener)
	ReadRepository() error
	IsCleanRepo() (bool, error)
	AddAll() error
	Commit(message string, amend bool) error
	CreateTag(tag string, message string) error
	GetLatestTag() (string, error)
	GetHeadCommit() (string, error)
	Push() error
	GetCommitsMessagesFromTagToHead(tag string) ([]string, error)
	GetName() string
	GetEmail() string
}

var g IVcs

type GitImpl struct {
	name             string
	email            string
	repository       *git.Repository
	repositoryOpener IRepositoryOpener
}

func init() {
	g = build()
}
func (g GitImpl) SetRepositoryOpener(repositoryOpener IRepositoryOpener) {
	g.repositoryOpener = repositoryOpener
}
func (g GitImpl) ReadRepository() error {

	err := g.repositoryOpener.Open(g)
	if err != nil {
		// RepositoryNotExists
		return err
	}

	cfg, err := g.repository.ConfigScoped(config.GlobalScope)
	if err != nil {
		return err
	}

	user := cfg.User
	if user.Name == "" && user.Email == "" {
		return nil
	}

	g.name = user.Name
	g.email = user.Email

	return nil
}
func (g GitImpl) GetStatus() (git.Status, error) {
	w, err := g.repository.Worktree()
	if err != nil {
		return nil, err
	}
	return w.Status()
}
func (g GitImpl) IsCleanRepo() (bool, error) {
	status, err := g.GetStatus()
	if err != nil {
		return false, err
	}
	return status.IsClean(), nil
}
func (g GitImpl) AddAll() error {
	w, err := g.repository.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(".")
	if err != nil {
		return err
	}
	return nil
}
func (g GitImpl) Commit(message string, amend bool) error {
	commitOptions := &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.name,
			Email: g.email,
			When:  time.Now(),
		},
		Amend: amend,
	}

	w, err := g.repository.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Commit(message, commitOptions)
	if err != nil {
		return err
	}
	return nil
}
func (g GitImpl) CreateTag(tag string, message string) error {

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
func (g GitImpl) GetLatestTag() (string, error) {

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
			// Es könnte ein leichtes Tag sein, also versuchen Sie, den Commit direkt zu bekommen
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
func (g GitImpl) Push() error {

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
func (g GitImpl) GetHeadCommit() (string, error) {
	headRef, err := g.repository.Head()
	if err != nil {
		return "", err
	}

	commit, err := g.repository.CommitObject(headRef.Hash())
	if err != nil {

	}

	return commit.Message, nil
}
func (g GitImpl) GetTag(tag string) (plumbing.Hash, error) {
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
func (g GitImpl) GetCommitsMessagesFromTagToHead(tag string) ([]string, error) {

	headRef, err := g.repository.Head()
	if err != nil {
		return nil, err
	}
	startHash := headRef.Hash()

	endHash, _ := g.GetTag(tag)
	if err != nil {
		endHash = plumbing.Hash{}
	}

	logOptions := &git.LogOptions{From: startHash}
	commitIter, err := g.repository.Log(logOptions)
	if err != nil {
		return nil, err
	}

	var commitMessages []string
	err = commitIter.ForEach(func(c *object.Commit) error {
		commitMessages = append(commitMessages, c.Message)
		if c.Hash == endHash {
			return fmt.Errorf("BREAK")
		}
		return nil
	})

	if err != nil && err.Error() != "BREAK" {
		return nil, err
	}

	return commitMessages, nil
}
func (g GitImpl) GetName() string {
	return g.name
}
func (g GitImpl) GetEmail() string {
	return g.email
}
func Get() IVcs {
	return g
}
func build() IVcs {
	g := new(GitImpl)
	return g
}
