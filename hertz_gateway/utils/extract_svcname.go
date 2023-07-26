package utils

import "path/filepath"

func ExtractServiceName(idlPath string) string {
	fileName := filepath.Base(idlPath)
	svcName := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	return svcName
}
