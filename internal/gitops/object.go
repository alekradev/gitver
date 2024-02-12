package gitops

import (
	"github.com/go-git/go-git/v5"
)

type GitRepositoryOpener struct {
	path string
}

func (o GitRepositoryOpener) SetPath(path string) {
	o.path = path
}

func (o GitRepositoryOpener) Open(vcs IVcs) error {
	repository, err := git.PlainOpen(o.path)
	if err != nil {
		return err
	}
	if g, ok := vcs.(GitImpl); ok {
		g.repository = repository
		return nil
	}

	return nil

}
