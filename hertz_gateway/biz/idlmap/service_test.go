package idlmap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/test"
	mymock "github.com/Linda-ui/orbital_HeBao/hertz_gateway/test/mock"
	"github.com/stretchr/testify/mock"
)

func TestManager_AddAllServices(t *testing.T) {
	// Create a temp directory with a subdirectory for testing
	tempDir := t.TempDir()
	subDir := tempDir + "/subdir"
	err := os.MkdirAll(subDir, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	files := []test.TempFile{
		{Name: "file1.thrift", Path: filepath.Join(tempDir, "file1.thrift"), Content: []byte(``)},
		{Name: "file2.thrift", Path: filepath.Join(tempDir, "file2.thrift"), Content: []byte(``)},
		{Name: "file3.thrift", Path: filepath.Join(subDir, "file3.thrift"), Content: []byte(``)},
		{Name: "file4.thrift", Path: filepath.Join(subDir, "file4.thrift"), Content: []byte(``)},
		{Name: "file5.thrift", Path: filepath.Join(tempDir, "file5.thrift"), Content: []byte(``)},
	}
	test.CreateTestFiles(t, files)

	mockRepo := mymock.NewRepository()

	// Use a MockMap to assert expectations by AddAll(...)
	testManager := manager{
		repo: mockRepo,
	}

	mockRepo.On("AddService", files[0].Path, mock.Anything).Return(nil).Once()
	mockRepo.On("AddService", files[1].Path, mock.Anything).Return(nil).Once()
	mockRepo.On("AddService", files[2].Path, mock.Anything).Return(nil).Once()
	mockRepo.On("AddService", files[3].Path, mock.Anything).Return(nil).Once()
	mockRepo.On("AddService", files[4].Path, mock.Anything).Return(nil).Once()

	testManager.AddAllServices(tempDir)

	mockRepo.AssertExpectations(t)
}
