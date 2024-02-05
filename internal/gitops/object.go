package gitops

import "github.com/go-git/go-git/v5"

type GitObject struct{}

func (o *GitObject) Open(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}
