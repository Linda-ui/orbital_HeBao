package utils

import (
	"path"
	"runtime"
)

func GetProjectIDLRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "./idl"
	}

	projectIDLRootPath := path.Join(path.Dir(filename), "..", "..", "idl")
	return projectIDLRootPath
}
