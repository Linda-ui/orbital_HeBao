package idlmap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/test"
	mymock "github.com/Linda-ui/orbital_HeBao/hertz_gateway/test/mock"
	"github.com/stretchr/testify/assert"
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

	testManager.AddAllServices(tempDir)

	mockRepo.AssertExpectations(t)
}

func TestManager_DynamicUpdate(t *testing.T) {
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
	}
	test.CreateTestFiles(t, files)

	mockRepo := mymock.NewRepository()

	// Use a MockMap to assert expectations by AddAll(...)
	testManager := manager{
		repo: mockRepo,
	}

	testManager.DynamicUpdate(tempDir)

	type args struct {
		idlPath string
	}

	tests := []struct {
		Name      string
		args      args
		svcName   string
		wantError bool
	}{
		{
			Name:      "valid thrift IDL file",
			args:      args{file1.Path},
			svcName:   "file1",
			wantError: false,
		},
		{
			Name:      "invalid thrift IDL file",
			args:      args{file2.Path},
			svcName:   "file2",
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := db.AddService(test.args.idlPath)
			if gotError := (err != nil); test.wantError != gotError {
				t.Errorf("got error = %v, want error = %v", gotError, test.wantError)
			}
			if !test.wantError {
				_, ok := db[test.svcName]
				assert.True(t, ok)
			}
		})
	}
}
