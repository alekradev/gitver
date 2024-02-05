package version

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

var testFilePath = "/etc/.gitver"
var testFileName = ".version.yaml"

var validTestFile = []byte(`
Version: 2024-02-04
History:
  - Version:
      Major: 0
      Minor: 0
      Patch: 0
    Date: 2024-02-04
  - Version:
      Major: 1
      Minor: 0
      Patch: 0
    Date: 2024-02-04
`)

var invalidTestFile = []byte(`
invalid content
`)

// AbsFilePath calls filepath.Abs on path.
func AbsFilePath(t *testing.T, path string) string {
	t.Helper()

	s, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}

	return s
}

func TestNewData(t *testing.T) {
	d := New()
	assert.NotNil(t, d)
	assert.Equal(t, "2024-02-04", d.Version)
	assert.Len(t, d.History, 1)
	assert.Equal(t, 0, d.History[0].Version.Major)
	// Füge weitere Assertions hinzu, um die Initialwerte zu überprüfen
}

func TestReadVersionSuccess(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	file.Close()

	v = New()

	v.setFileSystem(fs)
	v.setFilePath(testFilePath)
	v.setFileName(testFileName)

	err = v.readFile()
	assert.NoError(t, err)
	assert.Len(t, v.History, 2)
	assert.Equal(t, 0, v.History[0].Version.Major)
	assert.Equal(t, 0, v.History[0].Version.Minor)
	assert.Equal(t, 0, v.History[0].Version.Patch)
	assert.Equal(t, 1, v.History[1].Version.Major)
	assert.Equal(t, 0, v.History[1].Version.Minor)
	assert.Equal(t, 0, v.History[1].Version.Patch)
}

func TestReadVersionFileNotFound(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	v = New()

	v.setFileSystem(fs)
	v.setFilePath(testFilePath)
	v.setFileName(testFileName)

	err = v.readFile()
	assert.Error(t, err)
	// Überprüfe, ob der Fehler FileNotFoundError ist
}

func TestReadVersionInvalidFormat(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
	require.NoError(t, err)

	_, err = file.Write(invalidTestFile)
	require.NoError(t, err)

	file.Close()

	v = New()

	v.setFileSystem(fs)
	v.setFilePath(testFilePath)
	v.setFileName(testFileName)

	err = v.readFile()
	assert.Error(t, err)
	// Überprüfe, ob der Fehler FileFormatError ist
}

func TestWriteVersionSuccess(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	v = New()

	v.setFileSystem(fs)
	v.setFilePath(testFilePath)
	v.setFileName(testFileName)

	err = v.writeFile()
	assert.NoError(t, err)

	filePath := filepath.Join(testFilePath, testFileName)
	exists, _ := afero.Exists(v.fs, filePath)
	assert.True(t, exists)

	err = v.readFile()
	require.NoError(t, err)
	assert.Len(t, v.History, 1)
	assert.Equal(t, 0, v.History[0].Version.Major)
	assert.Equal(t, 0, v.History[0].Version.Minor)
	assert.Equal(t, 0, v.History[0].Version.Patch)
}

func TestWriteVersionCreateFolder(t *testing.T) {
	fs := afero.NewMemMapFs()

	v = New()

	v.setFileSystem(fs)
	v.setFilePath(testFilePath)
	v.setFileName(testFileName)

	err := v.writeFile()
	assert.NoError(t, err)

	filePath := filepath.Join(testFilePath, testFileName)
	exists, _ := afero.Exists(v.fs, filePath)
	assert.True(t, exists)

	err = v.readFile()
	require.NoError(t, err)
	assert.Len(t, v.History, 1)
	assert.Equal(t, 0, v.History[0].Version.Major)
	assert.Equal(t, 0, v.History[0].Version.Minor)
	assert.Equal(t, 0, v.History[0].Version.Patch)
}

func TestSafeWriteFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	file.Close()

	v = New()

	v.setFileSystem(fs)
	v.setFilePath(testFilePath)
	v.setFileName(testFileName)

	err = v.safeWriteFile()
	assert.Error(t, err)
}

func TestGetVersion(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	file.Close()

	v = New()

	v.setFileSystem(fs)
	v.setFilePath(testFilePath)
	v.setFileName(testFileName)

	err = v.readFile()

	version := v.getVersion()
	assert.Equal(t, "1.0.0", version)
}

func TestGetLastVersion(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	file.Close()

	v = New()

	v.setFileSystem(fs)
	v.setFilePath(testFilePath)
	v.setFileName(testFileName)

	err = v.readFile()

	lastVersion := v.getLastVersion()
	assert.Equal(t, "0.0.0", lastVersion)
}

func TestBumpMajor(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	file.Close()

	v = New()

	v.setFileSystem(fs)
	v.setFilePath(testFilePath)
	v.setFileName(testFileName)

	err = v.readFile()
	require.NoError(t, err)

	err = v.bumpMajor()
	assert.NoError(t, err)
	assert.Len(t, v.History, 3)
	assert.Equal(t, v.getVersion(), "2.0.0")
	assert.Equal(t, v.getLastVersion(), "1.0.0")
}

func TestBumpMinor(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	file.Close()

	v = New()

	v.setFileSystem(fs)
	v.setFilePath(testFilePath)
	v.setFileName(testFileName)

	err = v.readFile()
	require.NoError(t, err)

	err = v.bumpMinor()
	assert.NoError(t, err)
	assert.Len(t, v.History, 3)
	assert.Equal(t, v.getVersion(), "1.1.0")
	assert.Equal(t, v.getLastVersion(), "1.0.0")
}

func TestBumpPatch(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	file.Close()

	v = New()

	v.setFileSystem(fs)
	v.setFilePath(testFilePath)
	v.setFileName(testFileName)

	err = v.readFile()
	require.NoError(t, err)

	err = v.bumpPatch()
	assert.NoError(t, err)
	assert.Len(t, v.History, 3)
	assert.Equal(t, v.getVersion(), "1.0.1")
	assert.Equal(t, v.getLastVersion(), "1.0.0")
}

func TestSetDefaultVersion(t *testing.T) {
	v = New()

	err := v.setDefaultVersion("1.0.0")
	assert.NoError(t, err)
	assert.Len(t, v.History, 1)
	assert.Equal(t, v.getVersion(), "1.0.0")

}
