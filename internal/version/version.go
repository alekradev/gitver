package version

import (
	"fmt"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"time"
)

var v *Data

type Data struct {
	Version           string         `yaml:"Version"`
	History           []HistoryEntry `yaml:"History"`
	filepath          string
	filename          string
	fs                afero.Fs
	configPermissions os.FileMode
}

type HistoryEntry struct {
	Version Version `yaml:"Version"`
	Date    string  `yaml:"Date"`
}

type Version struct {
	Major int `yaml:"Major"`
	Minor int `yaml:"Minor"`
	Patch int `yaml:"Patch"`
}

// Private Functions
func init() {
	v = New()
}

func (v *Version) copy() Version {
	result := Version{}
	result.Major = v.Major
	result.Minor = v.Minor
	result.Patch = v.Patch
	return result
}

func (v *Version) toString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
func (d *Data) readFile() error {
	path := filepath.Join(d.filepath, d.filename)

	data, err := afero.ReadFile(d.fs, path)
	if err != nil {
		return FileNotFoundError(path)
	}

	err = yaml.Unmarshal(data, d)
	if err != nil {
		return FileFormatError(path)
	}

	return nil
}
func (d *Data) setFilePath(filepath string) {
	d.filepath = filepath
}
func (d *Data) setFileName(filename string) {
	d.filename = filename
}
func (d *Data) safeWriteFile() error {
	path := filepath.Join(d.filepath, d.filename)
	alreadyExists, err := afero.Exists(v.fs, path)
	if alreadyExists && err == nil {
		return FileAlreadyExistsError(path)
	}

	if err := d.writeFile(); err != nil {
		return err
	}

	return nil
}
func (d *Data) writeFile() error {
	path := filepath.Join(d.filepath, d.filename)

	data, err := yaml.Marshal(d)
	if err != nil {
		log.Fatalf("Fehler beim Konvertieren in YAML: %d", err)
		return err
	}

	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY

	f, err := d.fs.OpenFile(path, flags, d.configPermissions)
	if err != nil {
		return err
	}
	defer f.Close()

	err = afero.WriteFile(v.fs, path, data, os.ModePerm)
	if err != nil {
		return WriteOperationFailedError{path, err}
	}

	return nil
}
func (d *Data) bumpMajor() error {
	if len(d.History) == 0 {

	}

	version := d.History[len(d.History)-1].Version.copy()
	version.Major++
	version.Minor = 0
	version.Patch = 0
	entry := HistoryEntry{
		Version: version,
		Date:    time.Now().Format("YYYY-MM-DD"),
	}
	d.History = append(d.History, entry)

	return d.writeFile()
}
func (d *Data) bumpMinor() error {
	if len(d.History) == 0 {

	}

	version := d.History[len(d.History)-1].Version.copy()
	version.Minor++
	version.Patch = 0
	entry := HistoryEntry{
		Version: version,
		Date:    time.Now().Format("YYYY-MM-DD"),
	}
	d.History = append(d.History, entry)

	return d.writeFile()
}
func (d *Data) bumpPatch() error {
	if len(d.History) == 0 {

	}

	version := d.History[len(d.History)-1].Version.copy()
	version.Patch++
	entry := HistoryEntry{
		Version: version,
		Date:    time.Now().Format("YYYY-MM-DD"),
	}
	d.History = append(d.History, entry)

	return d.writeFile()
}
func (d *Data) setFileSystem(fs afero.Fs) {
	d.fs = fs
}
func (d *Data) getVersion() string {
	return d.History[len(d.History)-1].Version.toString()
}
func (d *Data) getLastVersion() string {
	return d.History[len(d.History)-2].Version.toString()
}
func (d *Data) setDefaultVersion(version string) error {
	_, err := fmt.Sscanf(version, "%d.%d.%d", &d.History[0].Version.Major, &d.History[0].Version.Minor, &d.History[0].Version.Patch)
	return err
}

// Public Functions

func New() *Data {
	v := new(Data)
	v.Version = "2024-02-04"
	v.History = []HistoryEntry{
		{
			Version: Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			Date: time.Now().Format("YYYY-MM-DD"),
		},
	}
	v.configPermissions = os.FileMode(0o644)
	v.fs = afero.NewOsFs()
	return v
}
func ReadFile() error {
	return v.readFile()
}
func SetFilePath(filePath string) {
	v.setFilePath(filePath)
}
func SetFileName(fileName string) {
	v.setFileName(fileName)
}
func SetFileSystem(fs afero.Fs) {
	v.setFileSystem(fs)
}
func SafeWriteFile() error {
	return v.safeWriteFile()
}
func BumpMajor() error {
	return v.bumpMajor()
}
func BumpMinor() error {
	return v.bumpMinor()
}
func BumpPatch() error {
	return v.bumpPatch()
}
func GetVersion() string {
	return v.getVersion()
}
func GetLastVersion() string {
	return v.getLastVersion()
}
func SetDefaultVersion(version string) error {
	return v.setDefaultVersion(version)
}
