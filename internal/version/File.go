package version

import (
	"fmt"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"time"
)

type IFile interface {
	ReadFile() error
	SetFilePath(filepath string)
	SetFileName(filename string)
	SetDefaultVersion(version string) error
	SafeWriteFile() error
	WriteFile() error
	BumpMajor() error
	BumpMinor() error
	BumpPatch() error
	SetFileSystem(fs afero.Fs)
	GetVersion() string
	GetLastVersion() string

	GetHistory() []IHistoryEntry
	GetFileSystem() afero.Fs
}

type File struct {
	Version           string          `yaml:"Version"`
	History           []IHistoryEntry `yaml:"History"`
	filepath          string
	filename          string
	fs                afero.Fs
	configPermissions os.FileMode
}

func (f *File) GetFileSystem() afero.Fs {
	return f.fs
}

func (f *File) GetHistory() []IHistoryEntry {
	return f.History
}
func (f *File) ReadFile() error {
	path := filepath.Join(f.filepath, f.filename)

	data, err := afero.ReadFile(f.fs, path)
	if err != nil {
		return FileNotFoundError(path)
	}

	err = yaml.Unmarshal(data, f)
	if err != nil {
		return FileFormatError(path)
	}

	return nil
}
func (f *File) SetFilePath(filepath string) {
	f.filepath = filepath
}
func (f *File) SetFileName(filename string) {
	f.filename = filename
}
func (f *File) SafeWriteFile() error {
	path := filepath.Join(f.filepath, f.filename)
	alreadyExists, err := afero.Exists(f.fs, path)
	if alreadyExists && err == nil {
		return FileAlreadyExistsError(path)
	}

	if err := f.WriteFile(); err != nil {
		return err
	}

	return nil
}
func (f *File) WriteFile() error {
	path := filepath.Join(f.filepath, f.filename)

	data, err := yaml.Marshal(f)
	if err != nil {
		log.Fatalf("Fehler beim Konvertieren in YAML: %d", err)
		return err
	}

	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY

	file, err := f.fs.OpenFile(path, flags, f.configPermissions)
	if err != nil {
		return err
	}
	defer file.Close()

	err = afero.WriteFile(f.fs, path, data, os.ModePerm)
	if err != nil {
		return WriteOperationFailedError{path, err}
	}

	return nil
}
func (f *File) BumpMajor() error {
	if len(f.History) == 0 {

	}

	version := f.History[len(f.History)-1].GetVersion().Copy()
	version.SetMajor(version.GetMajor() + 1)
	version.SetMinor(0)
	version.SetPatch(0)
	entry := HistoryEntry{
		Version: version,
		Date:    time.Now().Format("YYYY-MM-DD"),
	}
	f.History = append(f.History, &entry)

	return f.WriteFile()
}
func (f *File) BumpMinor() error {
	if len(f.History) == 0 {

	}

	version := f.History[len(f.History)-1].GetVersion().Copy()
	version.SetMinor(version.GetMinor() + 1)
	version.SetPatch(0)
	entry := HistoryEntry{
		Version: version,
		Date:    time.Now().Format("YYYY-MM-DD"),
	}
	f.History = append(f.History, &entry)

	return f.WriteFile()
}
func (f *File) BumpPatch() error {
	if len(f.History) == 0 {

	}

	version := f.History[len(f.History)-1].GetVersion().Copy()
	version.SetPatch(version.GetPatch() + 1)
	entry := HistoryEntry{
		Version: version,
		Date:    time.Now().Format("YYYY-MM-DD"),
	}
	f.History = append(f.History, &entry)

	return f.WriteFile()
}
func (f *File) SetFileSystem(fs afero.Fs) {
	f.fs = fs
}
func (f *File) GetVersion() string {
	return f.History[len(f.History)-1].GetVersion().ToString()
}
func (f *File) GetLastVersion() string {
	return f.History[len(f.History)-2].GetVersion().ToString()
}
func (f *File) SetDefaultVersion(version string) error {
	major := 0
	minor := 0
	patch := 0
	_, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	f.GetHistory()[0].GetVersion().SetMajor(major)
	f.GetHistory()[0].GetVersion().SetMinor(minor)
	f.GetHistory()[0].GetVersion().SetPatch(patch)
	return err
}

func Get() IFile {
	if f == nil {
		f = build()
	}
	return f
}
func build() IFile {
	v := new(Version)

	h := HistoryEntryBuilder().WithVersion(v).WithDate(time.Now().Format("YYYY-MM-DD")).Build()
	f := new(File)

	f.Version = "2024-02-04"
	f.History = []IHistoryEntry{
		h,
	}
	f.configPermissions = os.FileMode(0o644)
	f.fs = afero.NewOsFs()
	return f
}
