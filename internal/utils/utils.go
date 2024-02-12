package utils

import (
	"runtime"
)

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
