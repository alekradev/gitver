package gitops

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type Object interface {
	Open(path string) (*git.Repository, error)
}

type Repository interface {
	Worktree() (*git.Worktree, error)
	Head() (*plumbing.Reference, error)
	CommitObject(hash plumbing.Hash) (*object.Commit, error)
	Tags() (storer.ReferenceIter, error)
	TagObject(hash plumbing.Hash) (*object.Tag, error)
	CreateTag(name string, hash plumbing.Hash, opts *git.CreateTagOptions) (*plumbing.Reference, error)
	Push(options *git.PushOptions) error
	ConfigScoped(scope config.Scope) (*config.Config, error)
	Reference(name plumbing.ReferenceName, resolved bool) (*plumbing.Reference, error)
	Tag(name string) (*plumbing.Reference, error)
	Log(o *git.LogOptions) (object.CommitIter, error)
}

type Worktree interface {
	Add(pattern string) (plumbing.Hash, error)
	Commit(message string, opts *git.CommitOptions) (plumbing.Hash, error)
	Status() (git.Status, error)
}

type ReferenceIter interface {
	Next() (*plumbing.Reference, error)
	ForEach(func(*plumbing.Reference) error) error
	Close()
}

type TagObject interface {
	Tagger() *object.Signature
}

type CommitObject interface {
	Committer() *object.Signature
	Hash() plumbing.Hash
}
