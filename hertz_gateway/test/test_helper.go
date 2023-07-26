package test

import (
	"path"
	"runtime"

	"github.com/cloudwego/hertz/pkg/common/errors"
)

func GetIDLRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.ErrNothingRead
	}

	idlRootPath := path.Join(path.Dir(filename), "testfiles")
	return idlRootPath, nil
}
