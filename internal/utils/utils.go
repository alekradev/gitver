package utils

import (
	"runtime"
)

type Config struct {
	Meta  MetaConfig  `yaml:"Meta"`
	Files FilesConfig `yaml:"Files"`
}

type MetaConfig struct {
	SetPoms bool `yaml:"SetPoms"`
}

type FilesConfig struct {
	Poms []PomConfig `yaml:"Poms"`
}

type PomConfig struct {
	File string `yaml:"File"`
	Path string `yaml:"Path"`
}

func GetCurrentFunctionName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "Unbekannt"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "Unbekannt"
	}

	return fn.Name()
}
