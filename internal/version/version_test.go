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
	d := build()
	assert.NotNil(t, d)
	assert.Equal(t, "2024-02-04", d.GetVersion())
	assert.Len(t, d.GetHistory(), 1)
	assert.Equal(t, 0, d.GetHistory()[0].GetVersion().GetMajor())
	assert.Equal(t, 0, d.GetHistory()[0].GetVersion().GetMinor())
	assert.Equal(t, 0, d.GetHistory()[0].GetVersion().GetPatch())
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

	err = file.Close()
	assert.NoError(t, err)

	f = build()

	f.SetFileSystem(fs)
	f.SetFilePath(testFilePath)
	f.SetFileName(testFileName)

	err = f.ReadFile()
	assert.NoError(t, err)
	assert.Len(t, f.GetHistory(), 2)
	assert.Equal(t, 0, f.GetHistory()[0].GetVersion().GetMajor())
	assert.Equal(t, 0, f.GetHistory()[0].GetVersion().GetMinor())
	assert.Equal(t, 0, f.GetHistory()[0].GetVersion().GetPatch())
	assert.Equal(t, 0, f.GetHistory()[1].GetVersion().GetMajor())
	assert.Equal(t, 0, f.GetHistory()[1].GetVersion().GetMinor())
	assert.Equal(t, 0, f.GetHistory()[1].GetVersion().GetPatch())
}

func TestReadVersionFileNotFound(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	f = build()

	f.SetFileSystem(fs)
	f.SetFilePath(testFilePath)
	f.SetFileName(testFileName)

	err = f.ReadFile()
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

	err = file.Close()
	assert.NoError(t, err)

	f = build()

	f.SetFileSystem(fs)
	f.SetFilePath(testFilePath)
	f.SetFileName(testFileName)

	err = f.ReadFile()
	assert.Error(t, err)
	// Überprüfe, ob der Fehler FileFormatError ist
}

func TestWriteVersionSuccess(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	f = build()

	f.SetFileSystem(fs)
	f.SetFilePath(testFilePath)
	f.SetFileName(testFileName)

	err = f.WriteFile()
	assert.NoError(t, err)

	filePath := filepath.Join(testFilePath, testFileName)
	exists, _ := afero.Exists(f.GetFileSystem(), filePath)
	assert.True(t, exists)

	err = f.ReadFile()
	require.NoError(t, err)
	assert.Len(t, f.GetHistory(), 1)
	assert.Equal(t, 0, f.GetHistory()[0].GetVersion().GetMajor())
	assert.Equal(t, 0, f.GetHistory()[0].GetVersion().GetMinor())
	assert.Equal(t, 0, f.GetHistory()[0].GetVersion().GetPatch())
}

func TestWriteVersionCreateFolder(t *testing.T) {
	fs := afero.NewMemMapFs()

	f = build()

	f.SetFileSystem(fs)
	f.SetFilePath(testFilePath)
	f.SetFileName(testFileName)

	err := f.WriteFile()
	assert.NoError(t, err)

	filePath := filepath.Join(testFilePath, testFileName)
	exists, _ := afero.Exists(f.GetFileSystem(), filePath)
	assert.True(t, exists)

	err = f.ReadFile()
	require.NoError(t, err)
	assert.Len(t, f.GetHistory(), 1)
	assert.Equal(t, 0, f.GetHistory()[0].GetVersion().GetMajor())
	assert.Equal(t, 0, f.GetHistory()[0].GetVersion().GetMinor())
	assert.Equal(t, 0, f.GetHistory()[0].GetVersion().GetPatch())
}

func TestSafeWriteFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	f = build()

	f.SetFileSystem(fs)
	f.SetFilePath(testFilePath)
	f.SetFileName(testFileName)

	err = f.SafeWriteFile()
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

	err = file.Close()
	assert.NoError(t, err)

	f = build()

	f.SetFileSystem(fs)
	f.SetFilePath(testFilePath)
	f.SetFileName(testFileName)

	err = f.ReadFile()

	version := f.GetVersion()
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

	err = file.Close()
	assert.NoError(t, err)

	f = build()

	f.SetFileSystem(fs)
	f.SetFilePath(testFilePath)
	f.SetFileName(testFileName)

	err = f.ReadFile()

	lastVersion := f.GetLastVersion()
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

	err = file.Close()
	assert.NoError(t, err)

	f = build()

	f.SetFileSystem(fs)
	f.SetFilePath(testFilePath)
	f.SetFileName(testFileName)

	err = f.ReadFile()
	require.NoError(t, err)

	err = f.BumpMajor()
	assert.NoError(t, err)
	assert.Len(t, f.GetHistory(), 3)
	assert.Equal(t, f.GetVersion(), "2.0.0")
	assert.Equal(t, f.GetLastVersion(), "1.0.0")
}

func TestBumpMinor(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	f = build()

	f.SetFileSystem(fs)
	f.SetFilePath(testFilePath)
	f.SetFileName(testFileName)

	err = f.ReadFile()
	require.NoError(t, err)

	err = f.BumpMinor()
	assert.NoError(t, err)
	assert.Len(t, f.GetHistory(), 3)
	assert.Equal(t, f.GetVersion(), "1.1.0")
	assert.Equal(t, f.GetLastVersion(), "1.0.0")
}

func TestBumpPatch(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	f = build()

	f.SetFileSystem(fs)
	f.SetFilePath(testFilePath)
	f.SetFileName(testFileName)

	err = f.ReadFile()
	require.NoError(t, err)

	err = f.BumpPatch()
	assert.NoError(t, err)
	assert.Len(t, f.GetHistory(), 3)
	assert.Equal(t, f.GetVersion(), "1.0.1")
	assert.Equal(t, f.GetLastVersion(), "1.0.0")
}

func TestSetDefaultVersion(t *testing.T) {
	f = build()

	err := f.SetDefaultVersion("1.0.0")
	assert.NoError(t, err)
	assert.Len(t, f.GetHistory(), 1)
	assert.Equal(t, f.GetVersion(), "1.0.0")

}
