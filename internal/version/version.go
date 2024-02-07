package version

import "fmt"

var f IFile

type IVersion interface {
	Copy() IVersion
	ToString() string
	GetMajor() int
	GetMinor() int
	GetPatch() int
	SetMajor(value int)
	SetMinor(value int)
	SetPatch(value int)
}

type Version struct {
	Major int `yaml:"Major"`
	Minor int `yaml:"Minor"`
	Patch int `yaml:"Patch"`
}

// Private Functions
func init() {
	f = build()
}
func (v *Version) Copy() IVersion {
	result := new(Version)
	result.SetMajor(v.Major)
	result.SetMinor(v.Minor)
	result.SetPatch(v.Patch)
	return result
}

func (v *Version) ToString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) GetMajor() int {
	return v.Major
}
func (v *Version) GetMinor() int {
	return v.Minor
}
func (v *Version) GetPatch() int {
	return v.Patch
}
func (v *Version) SetMajor(value int) {
	v.Major = value
}
func (v *Version) SetMinor(value int) {
	v.Minor = value
}
func (v *Version) SetPatch(value int) {
	v.Patch = value
}
