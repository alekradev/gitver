package filesystem

import (
	"os"
	"path/filepath"
)

type Filesystem interface {
	Open(name string) (*os.File, error)
	Create(name string) (*os.File, error)
	Getwd() (dir string, err error)
	Stat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
	Mkdir(name string, perm os.FileMode) error
	IsNotExist(err error) bool
	WriteFile(name string, data []byte, perm os.FileMode) error
	Exists(path string) (bool, error)
	Find(dir string) (string, error)
}

type OsFilesystem struct{}

func New() Filesystem {
	return &OsFilesystem{}
}

func (fs OsFilesystem) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (fs OsFilesystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (fs OsFilesystem) Getwd() (string, error) {
	return os.Getwd()
}

func (fs OsFilesystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
func (fs OsFilesystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
func (fs OsFilesystem) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}
func (fs OsFilesystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
func (fs OsFilesystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (fs OsFilesystem) Exists(path string) (bool, error) {
	_, err := fs.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (fs OsFilesystem) Find(folderName string) (string, error) {
	dir, err := fs.Getwd()
	if err != nil {
		return "", UnhandledError(err.Error())
	}

	for {

		exist, err := fs.Exists(filepath.Join(dir, folderName))
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
