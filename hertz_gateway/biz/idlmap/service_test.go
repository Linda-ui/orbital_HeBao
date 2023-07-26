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

	mockRepo := mymock.NewRepository()

	// Use a MockMap to assert expectations by AddAll(...)
	testManager := manager{
		repo: mockRepo,
	}

	testIDLRoot, err := test.GetIDLRoot()
	if err != nil {
		t.Fatalf("failed to get IDL directory: %v", err)
	}

	// recursively find all files / directories
	err = filepath.Walk(testIDLRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("Error accessing path: %v\n", err)
		}

		if !info.Mode().IsDir() && info.Mode().IsRegular() {
			mockRepo.On("AddService", path, mock.Anything).Return(nil).Once()
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Error walking through directory: %v\n", err)
	}

	testManager.AddAllServices(testIDLRoot)

	mockRepo.AssertExpectations(t)
}
