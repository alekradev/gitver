package gitops

type IRepositoryOpener interface {
	Open(vcs IVcs) error
}
