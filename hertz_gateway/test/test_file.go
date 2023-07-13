package test

import (
	"os"
	"testing"
)

// define a TempFile struct type to store test TempFile information.
type TempFile struct {
	Name    string
	Path    string
	Content []byte
}

// a helper function for creating temporary test files.
func CreateTestFiles(t *testing.T, files []TempFile) {
	t.Helper()
	for _, file := range files {
		realFile, err := os.Create(file.Path)
		if err != nil {
			t.Fatalf("failed to create file '%s': %v", file.Path, err)
		}
		defer realFile.Close()

		_, err = realFile.Write(file.Content)
		if err != nil {
			t.Fatalf("failed to write file '%s': %v", file.Path, err)
		}
	}
}
