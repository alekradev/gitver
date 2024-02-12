package version

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

var testFilePath = "/etc/.gitver"
var testFileName = ".version"
var testFileExt = ".yaml"

var validTestFile = []byte(`
version: 2024-02-04
release:
  version: 1.0.0
config:
  default:
    version: 0.0.0
  commands:
    major: [bump major]
    minor: [bump minor]
    patch: [bump patch]
    auto: [bump auto]
  vcs:
    messages:
      commit: "Bump Version [%s] -> [%s]"
      tag: "Tagged by gitver"
    format:
      releaseTag: "r%s"
      versionTag: "v%s"
  files:
    - filepath: ./test/pom.xml
      format: xml
      xpath: /root/version
    - filepath: ./test/test.json
      format: json
      jsonpath:
    - filepath: ./test/test.txt
      format: text
      placeholder: {VERSION}
`)

var invalidTestFile = []byte(`
invalid content
`)

// AbsFilePath calls filepath.Abs on path.
//func AbsFilePath(t *testing.T, path string) string {
//	t.Helper()
//
//	s, err := filepath.Abs(path)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	return s
//}

func TestNewData(t *testing.T) {
	g := build()
	assert.NotNil(t, g)
	assert.Equal(t, "0.0.0", g.GetVersion())
}

func TestGitVer_GetVersion(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	path, err := filepath.Abs(filepath.Join(testFilePath, testFileName, testFileExt))
	require.NoError(t, err)

	file, err := fs.Create(path)
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	gitVer := Get()
	gitVer.SetConfigFile(path)
	gitVer.SetFileSystem(fs)

	err = gitVer.ReadFile()
	require.NoError(t, err)

	assert.Equal(t, "1.0.0", gitVer.GetVersion())
}

func TestVersion_FromString(t *testing.T) {

}

func TestGitVer_BumpMajor(t *testing.T) {

	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	path, err := filepath.Abs(filepath.Join(testFilePath, testFileName, testFileExt))
	require.NoError(t, err)

	file, err := fs.Create(path)
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	gitVer := Get()
	gitVer.SetConfigFile(path)
	gitVer.SetFileSystem(fs)

	err = gitVer.ReadFile()
	require.NoError(t, err)

	err = gitVer.BumpMajor()
	require.NoError(t, err)

	assert.Equal(t, "2.0.0", gitVer.GetVersion())
}

func TestGitVer_BumpMinor(t *testing.T) {

	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	path, err := filepath.Abs(filepath.Join(testFilePath, testFileName, testFileExt))
	require.NoError(t, err)

	file, err := fs.Create(path)
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	gitVer := Get()
	gitVer.SetConfigFile(path)
	gitVer.SetFileSystem(fs)

	err = gitVer.ReadFile()
	require.NoError(t, err)

	err = gitVer.BumpMinor()
	require.NoError(t, err)

	assert.Equal(t, "1.1.0", gitVer.GetVersion())
}

func TestGitVer_BumpPatch(t *testing.T) {

	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	path, err := filepath.Abs(filepath.Join(testFilePath, testFileName, testFileExt))
	require.NoError(t, err)

	file, err := fs.Create(path)
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	gitVer := Get()
	gitVer.SetConfigFile(path)
	gitVer.SetFileSystem(fs)

	err = gitVer.ReadFile()
	require.NoError(t, err)

	err = gitVer.BumpPatch()
	require.NoError(t, err)

	assert.Equal(t, "1.0.1", gitVer.GetVersion())
}

func TestGitVer_BumpAuto(t *testing.T) {

	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	path, err := filepath.Abs(filepath.Join(testFilePath, testFileName, testFileExt))
	require.NoError(t, err)

	file, err := fs.Create(path)
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	gitVer := Get()
	gitVer.SetConfigFile(path)
	gitVer.SetFileSystem(fs)

	err = gitVer.ReadFile()
	require.NoError(t, err)

	err = gitVer.BumpAuto()
	require.NoError(t, err)

	assert.Equal(t, "1.0.1", gitVer.GetVersion())
}

func TestGitVer_BumpFlex(t *testing.T) {

	fs := afero.NewMemMapFs()

	err := fs.Mkdir(testFilePath, 0o777)
	require.NoError(t, err)

	path, err := filepath.Abs(filepath.Join(testFilePath, testFileName, testFileExt))
	require.NoError(t, err)

	file, err := fs.Create(path)
	require.NoError(t, err)

	_, err = file.Write(validTestFile)
	require.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	gitVer := Get()
	gitVer.SetConfigFile(path)
	gitVer.SetFileSystem(fs)

	err = gitVer.ReadFile()
	require.NoError(t, err)

	err = gitVer.BumpFlex()
	require.NoError(t, err)

	assert.Equal(t, "1.0.1", gitVer.GetVersion())
}

//func TestReadVersionSuccess(t *testing.T) {
//	fs := afero.NewMemMapFs()
//	err := fs.Mkdir(testFilePath, 0o777)
//	require.NoError(t, err)
//	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
//	require.NoError(t, err)
//	_, err = file.Write(validTestFile)
//	require.NoError(t, err)
//	err = file.Close()
//	assert.NoError(t, err)
//
//	g = build()
//
//	g.SetFileSystem(fs)
//	g.SetFilePath(testFilePath)
//	g.SetFileName(testFileName)
//
//	err = g.ReadFile()
//	assert.NoError(t, err)
//
//}

//func TestReadVersionFileNotFound(t *testing.T) {
//	fs := afero.NewMemMapFs()
//
//	err := fs.Mkdir(testFilePath, 0o777)
//	require.NoError(t, err)
//
//	var g = build()
//
//	g.SetFileSystem(fs)
//	g.SetFilePath(testFilePath)
//	g.SetFileName(testFileName)
//
//	err = g.ReadFile()
//	assert.Error(t, err)
//}

//func TestReadVersionInvalidFormat(t *testing.T) {
//	fs := afero.NewMemMapFs()
//
//	err := fs.Mkdir(testFilePath, 0o777)
//	require.NoError(t, err)
//
//	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
//	require.NoError(t, err)
//
//	_, err = file.Write(invalidTestFile)
//	require.NoError(t, err)
//
//	err = file.Close()
//	assert.NoError(t, err)
//
//	var g = build()
//
//	g.SetFileSystem(fs)
//	g.SetFilePath(testFilePath)
//	g.SetFileName(testFileName)
//
//	err = g.ReadFile()
//	assert.Error(t, err)
//}

//func TestWriteVersionSuccess(t *testing.T) {
//	fs := afero.NewMemMapFs()
//
//	err := fs.Mkdir(testFilePath, 0o777)
//	require.NoError(t, err)
//
//	var g = build()
//
//	g.SetFileSystem(fs)
//	g.SetFilePath(testFilePath)
//	g.SetFileName(testFileName)
//
//	err = g.WriteFile()
//	assert.NoError(t, err)
//
//	filePath := filepath.Join(testFilePath, testFileName)
//	exists, _ := afero.Exists(fs, filePath)
//	assert.True(t, exists)
//
//	err = g.ReadFile()
//	require.NoError(t, err)
//}

//func TestWriteVersionCreateFolder(t *testing.T) {
//	fs := afero.NewMemMapFs()
//
//	var g = build()
//
//	g.SetFileSystem(fs)
//	g.SetFilePath(testFilePath)
//	g.SetFileName(testFileName)
//
//	err := g.WriteFile()
//	assert.NoError(t, err)
//
//	filePath := filepath.Join(testFilePath, testFileName)
//	exists, _ := afero.Exists(fs, filePath)
//	assert.True(t, exists)
//
//	err = g.ReadFile()
//	require.NoError(t, err)
//}

//func TestSafeWriteFile(t *testing.T) {
//	fs := afero.NewMemMapFs()
//
//	err := fs.Mkdir(testFilePath, 0o777)
//	require.NoError(t, err)
//
//	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
//	require.NoError(t, err)
//
//	_, err = file.Write(validTestFile)
//	require.NoError(t, err)
//
//	err = file.Close()
//	assert.NoError(t, err)
//
//	f = build()
//
//	f.SetFileSystem(fs)
//	f.SetFilePath(testFilePath)
//	f.SetFileName(testFileName)
//
//	err = f.SafeWriteFile()
//	assert.Error(t, err)
//}

//func TestGetVersion(t *testing.T) {
//	fs := afero.NewMemMapFs()
//
//	err := fs.Mkdir(testFilePath, 0o777)
//	require.NoError(t, err)
//
//	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
//	require.NoError(t, err)
//
//	_, err = file.Write(validTestFile)
//	require.NoError(t, err)
//
//	err = file.Close()
//	assert.NoError(t, err)
//
//	var g = build()
//
//	g.SetFileSystem(fs)
//	g.SetFilePath(testFilePath)
//	g.SetFileName(testFileName)
//
//	err = g.ReadFile()
//
//	version := g.GetVersion()
//	assert.Equal(t, "1.0.0", version)
//}
//
//func TestGetLastVersion(t *testing.T) {
//	fs := afero.NewMemMapFs()
//
//	err := fs.Mkdir(testFilePath, 0o777)
//	require.NoError(t, err)
//
//	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
//	require.NoError(t, err)
//
//	_, err = file.Write(validTestFile)
//	require.NoError(t, err)
//
//	err = file.Close()
//	assert.NoError(t, err)
//
//	var g = build()
//
//	g.SetFileSystem(fs)
//	g.SetFilePath(testFilePath)
//	g.SetFileName(testFileName)
//
//	err = g.ReadFile()
//
//	lastVersion := g.GetLastVersion()
//	assert.Equal(t, "0.0.0", lastVersion)
//}
//
//func TestBumpMajor(t *testing.T) {
//	fs := afero.NewMemMapFs()
//
//	err := fs.Mkdir(testFilePath, 0o777)
//	require.NoError(t, err)
//
//	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
//	require.NoError(t, err)
//
//	_, err = file.Write(validTestFile)
//	require.NoError(t, err)
//
//	err = file.Close()
//	assert.NoError(t, err)
//
//	f = build()
//
//	f.SetFileSystem(fs)
//	f.SetFilePath(testFilePath)
//	f.SetFileName(testFileName)
//
//	err = f.ReadFile()
//	require.NoError(t, err)
//
//	err = f.BumpMajor()
//	assert.NoError(t, err)
//	assert.Len(t, f.GetHistory(), 3)
//	assert.Equal(t, f.GetVersion(), "2.0.0")
//	assert.Equal(t, f.GetLastVersion(), "1.0.0")
//}
//
//func TestBumpMinor(t *testing.T) {
//	fs := afero.NewMemMapFs()
//
//	err := fs.Mkdir(testFilePath, 0o777)
//	require.NoError(t, err)
//
//	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
//	require.NoError(t, err)
//
//	_, err = file.Write(validTestFile)
//	require.NoError(t, err)
//
//	err = file.Close()
//	assert.NoError(t, err)
//
//	f = build()
//
//	f.SetFileSystem(fs)
//	f.SetFilePath(testFilePath)
//	f.SetFileName(testFileName)
//
//	err = f.ReadFile()
//	require.NoError(t, err)
//
//	err = f.BumpMinor()
//	assert.NoError(t, err)
//	assert.Len(t, f.GetHistory(), 3)
//	assert.Equal(t, f.GetVersion(), "1.1.0")
//	assert.Equal(t, f.GetLastVersion(), "1.0.0")
//}
//
//func TestBumpPatch(t *testing.T) {
//	fs := afero.NewMemMapFs()
//
//	err := fs.Mkdir(testFilePath, 0o777)
//	require.NoError(t, err)
//
//	file, err := fs.Create(AbsFilePath(t, filepath.Join(testFilePath, testFileName)))
//	require.NoError(t, err)
//
//	_, err = file.Write(validTestFile)
//	require.NoError(t, err)
//
//	err = file.Close()
//	assert.NoError(t, err)
//
//	f = build()
//
//	f.SetFileSystem(fs)
//	f.SetFilePath(testFilePath)
//	f.SetFileName(testFileName)
//
//	err = f.ReadFile()
//	require.NoError(t, err)
//
//	err = f.BumpPatch()
//	assert.NoError(t, err)
//	assert.Len(t, f.GetHistory(), 3)
//	assert.Equal(t, f.GetVersion(), "1.0.1")
//	assert.Equal(t, f.GetLastVersion(), "1.0.0")
//}
//
//func TestSetDefaultVersion(t *testing.T) {
//	f = build()
//
//	err := f.SetDefaultVersion("1.0.0")
//	assert.NoError(t, err)
//	assert.Len(t, f.GetHistory(), 1)
//	assert.Equal(t, f.GetVersion(), "1.0.0")
//
//}
