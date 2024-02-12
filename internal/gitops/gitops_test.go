package gitops

import (
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	username             = "User"
	email                = "user@example.org"
	exampleFileName      = "file%d.txt"
	exampleCommitMessage = "Test Commit"
)

var (
	commits = []string{
		"init",
		"feat: implement feature",
		"fix: issue",
		"fix: issue",
		"release: v1.1.0",
	}
	tags = []string{
		"r1.0.0",
		"v1.1.0",
	}
	tagMessages = []string{
		"release: r1.0.0",
		"version: v1.1.0",
	}
)

func createTestRepo() (*git.Repository, error) {
	// Initialisiere ein neues In-Memory-Repository
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	if err != nil {
		return nil, fmt.Errorf("failed to init repo: %w", err)
	}

	cfg := config.NewConfig()
	cfg.User.Name = username
	cfg.User.Email = email
	err = repo.Storer.SetConfig(cfg)

	// Erhalte den Worktree, um Commits hinzufügen zu können
	wt, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	// Erstelle 5 Commits
	var commitHashes []plumbing.Hash
	for i, commit := range commits {

		filename := fmt.Sprintf(exampleFileName, i)
		_, err = wt.Filesystem.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}

		// Füge die Datei zum Index hinzu
		_, err = wt.Add(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to add file to index: %w", err)
		}

		timeStamp := time.Now().Add(time.Duration(i) * time.Hour)
		hash, err := wt.Commit(commit, &git.CommitOptions{
			Author: &object.Signature{
				Name:  username,
				Email: email,
				When:  timeStamp,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to Commit: %w", err)
		}
		commitHashes = append(commitHashes, hash)
	}

	// Erstelle 2 Tags, eines für den ersten und eines für den letzten Commit
	_, err = repo.CreateTag(tags[0], commitHashes[1], &git.CreateTagOptions{Message: tagMessages[0]})
	if err != nil {
		return nil, fmt.Errorf("failed to create tag v1.0: %w", err)
	}
	_, err = repo.CreateTag(tags[1], commitHashes[3], &git.CreateTagOptions{Message: tagMessages[0]})
	if err != nil {
		return nil, fmt.Errorf("failed to create tag v2.0: %w", err)
	}

	return repo, nil
}

type MockRepositoryOpener struct {
	repository *git.Repository
}

func (o *MockRepositoryOpener) Open(vcs IVcs) error {
	repository, err := createTestRepo()
	if err != nil {
		return err
	}
	if g, ok := vcs.(GitImpl); ok {
		o.repository = repository
		g.repository = repository
		return nil
	}

	return nil
}

func TestReadRepository(t *testing.T) {
	obj := new(MockRepositoryOpener)
	g := build()
	g.SetRepositoryOpener(obj)
	err := g.ReadRepository()
	assert.NoError(t, err)
	assert.Equal(t, g.GetName(), username)
	assert.Equal(t, g.GetEmail(), email)
}

func TestIsCleanRepo(t *testing.T) {
	obj := new(MockRepositoryOpener)
	g := build()
	g.SetRepositoryOpener(obj)
	err := g.ReadRepository()
	require.NoError(t, err)

	isClean, err := g.IsCleanRepo()
	assert.NoError(t, err)
	assert.True(t, isClean)
}

func TestIsCleanRepoWithChanges(t *testing.T) {
	obj := new(MockRepositoryOpener)
	g := build()
	g.SetRepositoryOpener(obj)
	err := g.ReadRepository()
	require.NoError(t, err)
	w, err := obj.repository.Worktree()
	require.NoError(t, err)
	_, err = w.Filesystem.Create(exampleFileName)
	_, err = w.Add(exampleFileName)

	isClean, err := g.IsCleanRepo()
	assert.NoError(t, err)
	assert.False(t, isClean)
}

func TestAdd(t *testing.T) {
	obj := new(MockRepositoryOpener)
	g := build()
	g.SetRepositoryOpener(obj)
	err := g.ReadRepository()
	require.NoError(t, err)
	w, err := obj.repository.Worktree()
	require.NoError(t, err)
	_, err = w.Filesystem.Create(exampleFileName)
	require.NoError(t, err)

	err = g.AddAll()

	assert.NoError(t, err)
}

func TestCommit(t *testing.T) {
	obj := new(MockRepositoryOpener)
	g := build()
	g.SetRepositoryOpener(obj)
	err := g.ReadRepository()
	require.NoError(t, err)
	w, err := obj.repository.Worktree()
	require.NoError(t, err)
	_, err = w.Filesystem.Create(exampleFileName)
	_, err = w.Add(exampleFileName)

	err = g.Commit(exampleCommitMessage, false)
	isClean, err := g.IsCleanRepo()
	require.NoError(t, err)

	commits, err := g.GetCommitsMessagesFromTagToHead("EMPTY")
	require.NoError(t, err)

	assert.NoError(t, err)
	assert.True(t, isClean)
	assert.Len(t, commits, 6)
	assert.Equal(t, exampleCommitMessage, commits[0])
}

//func TestCommitAmend(t *testing.T) {
//	obj := new(MockRepositoryOpener)
//	g := build()
//	g.SetObject(obj)
//	err := g.ReadRepository()
//	require.NoError(t, err)
//
//	_, err = g.worktree.Filesystem.Create(exampleFileName)
//	require.NoError(t, err)
//
//	hash, err := g.worktree.Add(exampleFileName)
//	require.NoError(t, err)
//	require.NotEmpty(t, hash)
//
//	err = g.Commit(exampleCommitMessage, true)
//	isClean, err := g.IsCleanRepo()
//	require.NoError(t, err)
//
//	commits, err := g.GetCommitsBetweenTags(EMPTY, EMPTY)
//	require.NoError(t, err)
//
//	assert.NoError(t, err)
//	assert.True(t, isClean)
//	assert.Len(t, commits, 5)
//	assert.Equal(t, exampleCommitMessage, commits[0].Message)
//}

func TestCreateTag(t *testing.T) {
	obj := new(MockRepositoryOpener)
	g := build()
	g.SetRepositoryOpener(obj)
	err := g.ReadRepository()
	require.NoError(t, err)

	headRef, err := obj.repository.Head()
	require.NoError(t, err)

	tagName := "v2.0.0"
	tagMessage := "Initial release\n"

	err = g.CreateTag(tagName, tagMessage)
	assert.NoError(t, err)

	tagRef, err := obj.repository.Tag(tagName)
	require.NoError(t, err)

	tagObj, err := obj.repository.TagObject(tagRef.Hash())
	require.NoError(t, err)

	assert.Equal(t, tagMessage, tagObj.Message)
	assert.Equal(t, headRef.Hash(), tagObj.Target)
}

//func TestGetLatestTag(t *testing.T) {
//
//	obj := new(MockRepositoryOpener)
//
//	g := build()
//	g.SetObject(obj)
//	err := g.ReadRepository()
//	require.NoError(t, err)
//
//	// Führe die zu testende Funktion aus
//	tagName, err := g.GetLatestTag()
//
//	// Überprüfe das Ergebnis
//	assert.NoError(t, err)
//	assert.Equal(t, tags[1], tagName)
//}

//func TestGetTag(t *testing.T) {
//	obj := new(MockRepositoryOpener)
//	g := build()
//	g.SetRepositoryOpener(obj)
//	err := g.ReadRepository()
//	require.NoError(t, err)
//
//	hash, err := g.GetTag(tags[0])
//	assert.NoError(t, err)
//	assert.NotEmpty(t, hash)
//}

func TestGetCommitsBetweenTags(t *testing.T) {
	obj := new(MockRepositoryOpener)
	g := build()
	g.SetRepositoryOpener(obj)
	err := g.ReadRepository()
	require.NoError(t, err)

	commits, err := g.GetCommitsMessagesFromTagToHead(tags[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, commits)
	assert.Len(t, commits, 3)
}

//func TestGetCommitsBetweenFirstAndTag(t *testing.T) {
//	r, _ := createTestRepo()
//
//	g := build()
//	g.SetRepository(r)
//
//	commits, err := g.GetCommitsMessagesFromTagToHead(tags[1], EMPTY)
//	assert.NoError(t, err)
//	assert.NotEmpty(t, commits)
//	assert.Len(t, commits, 4)
//}

//func TestGetCommitsBetweenTagAndHead(t *testing.T) {
//	r, _ := createTestRepo()
//
//	g := build()
//	g.SetRepository(r)
//
//	commits, err := g.GetCommitsMessagesFromTagToHead(EMPTY, tags[0])
//	assert.NoError(t, err)
//	assert.NotEmpty(t, commits)
//	assert.Len(t, commits, 4)
//}

//func TestGetCommitsBetweenFirstAndHead(t *testing.T) {
//	r, _ := createTestRepo()
//
//	g := build()
//	g.SetRepository(r)
//
//	commits, err := g.GetCommitsMessagesFromTagToHead(EMPTY, EMPTY)
//	assert.NoError(t, err)
//	assert.NotEmpty(t, commits)
//	assert.Len(t, commits, 5)
//}
