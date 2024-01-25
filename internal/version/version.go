package version

import (
	"fmt"
	"github.com/spf13/afero"
	"gotver/internal/constants"
	"io"
	"log"
	"os"
	"path/filepath"
)

var v *Version

type Version struct {
	major               int
	minor               int
	patch               int
	versionFilePath     string
	versionFileName     string
	lastVersionFileName string
	lastVersion         string
	fs                  afero.Fs
}

func init() {
	v = New()
}

func New() *Version {
	v := new(Version)
	v.major = 0
	v.minor = 0
	v.patch = 0
	v.versionFilePath = ""
	v.versionFileName = ""
	v.lastVersionFileName = ".lastversion"
	v.lastVersion = "0.0.0"
	v.fs = afero.NewOsFs()
	return v
}

// GetProjectDirectory returns the directory containing the .gotver folder.
func GetProjectDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", UnhandledError(err.Error())
	}

	for {
		// Überprüfen Sie, ob das aktuelle Verzeichnis `.gotver` enthält.
		if _, err := os.Stat(filepath.Join(dir, constants.ConfigFolderName)); !os.IsNotExist(err) {
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

	return "", ProjectDirectoryNotFoundError(dir)
}

func ReadVersion() error {
	return v.ReadVersion()
}

func (v *Version) ReadVersion() error {
	versionFilePath := filepath.Join(v.versionFilePath, v.versionFileName)
	lastVersionFilePath := filepath.Join(v.versionFilePath, v.lastVersionFileName)

	data, err := os.ReadFile(versionFilePath)
	if err != nil {
		return FileNotFoundError(versionFilePath)
	}

	_, err = fmt.Sscanf(string(data), "%d.%d.%d", &v.major, &v.minor, &v.patch)
	if err != nil {
		return FileFormatError(versionFilePath)
	}

	if _, err := os.Stat(lastVersionFilePath); os.IsNotExist(err) {
		return nil
	}

	data, err = os.ReadFile(lastVersionFilePath)
	if err != nil {
		return FileNotFoundError(lastVersionFilePath)
	}

	_, err = fmt.Sscan(string(data), &v.lastVersion)
	if err != nil {
		return FileFormatError(versionFilePath)
	}

	return nil
}

func SetFilePath(filePath string) {
	v.SetFilePath(filePath)
}

func (v *Version) SetFilePath(filePath string) {
	v.versionFilePath = filePath
}

func SetFileName(fileName string) {
	v.SetFileName(fileName)
}

func (v *Version) SetFileName(fileName string) {
	v.versionFileName = fileName
}

func SafeWriteVersion() error {
	return v.SafeWriteVersion()
}
func (v *Version) SafeWriteVersion() error {
	dir := filepath.Join(v.versionFilePath)
	versionFilePath := filepath.Join(v.versionFilePath, v.versionFileName)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	alreadyExists, err := afero.Exists(v.fs, versionFilePath)
	if alreadyExists && err == nil {
		return FileAlreadyExistsError(versionFilePath)
	}

	if err := v.WriteVersion(); err != nil {
		return err
	}

	return nil
}

func ToString() string {
	return v.ToString()
}
func (v *Version) ToString() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

func GetLastVersion() string {
	return v.GetLastVersion()
}
func (v *Version) GetLastVersion() string {
	return v.lastVersion
}

func FromString(version string) error {
	return v.FromString(version)
}
func (v *Version) FromString(version string) error {
	_, err := fmt.Sscanf(version, "%d.%d.%d", &v.major, &v.minor, &v.patch)
	if err != nil {
		return InputValueError(version)
	}
	return nil
}

func WriteVersion() error {
	return v.WriteVersion()
}

func (v *Version) WriteVersion() error {
	versionFilePath := filepath.Join(v.versionFilePath, v.versionFileName)
	lastVersionFilePath := filepath.Join(v.versionFilePath, v.lastVersionFileName)

	if _, err := os.Stat(versionFilePath); !os.IsNotExist(err) {
		source, err := os.Open(versionFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer source.Close()

		// Zieldatei erstellen
		destination, err := os.Create(lastVersionFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer destination.Close()

		// Inhalt kopieren
		_, err = io.Copy(destination, source)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := os.WriteFile(versionFilePath, []byte(v.ToString()), os.ModePerm)
	if err != nil {
		return WriteOperationFailedError{versionFilePath, err}
	}

	return nil
}

func BumpMajor() error {
	return v.BumpMajor()
}
func (v *Version) BumpMajor() error {
	v.lastVersion = v.ToString()
	v.major++
	v.minor = 0
	v.patch = 0

	return v.WriteVersion()
}

func BumpMinor() error {
	return v.BumpMinor()
}
func (v *Version) BumpMinor() error {
	v.lastVersion = v.ToString()
	v.minor++
	v.patch = 0

	return v.WriteVersion()
}

func BumpPatch() error {
	return v.BumpPatch()
}
func (v *Version) BumpPatch() error {
	v.lastVersion = v.ToString()
	v.patch++

	return v.WriteVersion()
}
