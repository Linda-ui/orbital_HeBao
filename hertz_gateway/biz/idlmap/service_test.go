package idlmap

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

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

	root, err := test.GetIDLRoot()
	if err != nil {
		t.Fatalf("failed to get IDL directory: %v", err)
	}

	// recursively find all files / directories
	err = filepath.Walk(*root, func(path string, info os.FileInfo, err error) error {
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

	testManager.AddAllServices(*root)

	mockRepo.AssertExpectations(t)
}

func Test_manager_DynamicUpdate(t *testing.T) {
	// Initial set up
	mockRepo := mymock.NewRepository()
	testManager := manager{
		repo: mockRepo,
	}

	root, err := test.GetIDLRoot()
	if err != nil {
		t.Fatalf("failed to get IDL directory: %v", err)
	}

	go func() {
		testManager.DynamicUpdate(*root)
	}()

	// Add a new directory and add two new files in new directory
	// assert AddService called once for each file
	newDirPath := filepath.Join(*root, "new_test_directory")
	newFile1Path := filepath.Join(newDirPath, "new_file_1.thrift")
	newFile2Path := filepath.Join(newDirPath, "new_file_2.thrift")
	newFilePaths := []string{newFile1Path, newFile2Path, newFile2Path}
	newFileContent := ``

	err = os.Mkdir(newDirPath, os.FileMode(0777))
	if err != nil {
		t.Fatalf("Error creating the directory: %v", err)
		return
	}
	time.Sleep(10 * time.Millisecond)

	for _, newFilePath := range newFilePaths {
		mockRepo.On("AddService", newFilePath, mock.Anything).Return(nil).Once()
		err = ioutil.WriteFile(newFilePath, []byte(newFileContent), os.FileMode(0644))
		if err != nil {
			t.Fatalf("Error creating the file: %v", err)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	// delete a test file
	// assert deleteservice called once
	mockRepo.On("DeleteService", "new_file_2").Once()

	err = os.Remove(newFile2Path)
	if err != nil {
		t.Fatalf("Error deleting the file: %v", err)
		return
	}
	time.Sleep(10 * time.Millisecond)

	// delete a directory
	// assert deleteservice called once
	mockRepo.On("DeleteService", "new_file_1").Once()
	mockRepo.On("DeleteService", "new_test_directory").Once()

	err = os.RemoveAll(newDirPath)
	if err != nil {
		t.Fatalf("Error deleting the file: %v", err)
		return
	}
	time.Sleep(10 * time.Millisecond)

	// assert
	mockRepo.AssertExpectations(t)
}
