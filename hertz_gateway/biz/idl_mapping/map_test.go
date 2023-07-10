package idl_mapping

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient implements the genericclient.Client interface for testing.
type MockClient struct{}

func (c *MockClient) GenericCall(ctx context.Context, method string, request interface{}, callOptions ...callopt.Option) (response interface{}, err error) {
	return "test", nil
}

func (c *MockClient) Close() error {
	return nil
}

func TestDynamicMap_GetClient(t *testing.T) {
	dynamicMap := &DynamicMap{
		innerMap: make(map[string]genericclient.Client),
	}

	mockClient := &MockClient{}

	dynamicMap.innerMap["serviceA"] = mockClient

	client, ok := dynamicMap.GetClient("serviceA")
	assert.True(t, ok)
	assert.Equal(t, mockClient, client)

	client, ok = dynamicMap.GetClient("serviceB")
	assert.False(t, ok)
	assert.Nil(t, client)
}

// a helper method for checking whether the dynamic map contains the service.
func (m *DynamicMap) hasService(serviceName string) bool {
	_, ok := m.innerMap[serviceName]
	return ok
}

// a helper function for creating temporary test files.
func createTestFiles(t *testing.T, arr []TestCase) {
	for _, tc := range arr {
		createTestFile(t, tc.path, tc.name, tc.content)
	}
}

func createTestFile(t *testing.T, filePath string, fileName string, fileContent []byte) {

	fullFilePath := filepath.Join(filePath, fileName)

	file, err := os.Create(fullFilePath)
	if err != nil {
		t.Fatalf("failed to create file '%s': %v", fullFilePath, err)
	}

	defer file.Close()

	_, err = file.Write(fileContent)
	if err != nil {
		t.Fatalf("failed to write file '%s': %v", fullFilePath, err)
	}
}

// a test case struct for testing the Add and AddAll methods.
type TestCase struct {
	path    string
	name    string
	content []byte
}

func TestDynamicMap_Add(t *testing.T) {
	// Create a temp directory with temp files for testing.
	tempDir, err := os.MkdirTemp("", "dir")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	file1Path := filepath.Join(tempDir, "file1.thrift")
	file2Path := filepath.Join(tempDir, "file2.thrift")
	// create a new test table.
	tcArr := []TestCase{}

	// a valid thrift file with a service. Client is expected to be created.
	file1Content := []byte(`
		namespace go example.file1

		service MyService {
			i64 add(1: i64 a, 2: i64 b)
		}
	`)
	tcArr = append(tcArr, TestCase{tempDir, "file1.thrift", file1Content})

	// an invalid thrift file without a service. Error is expected to be returned.
	file2Content := []byte(``)
	tcArr = append(tcArr, TestCase{tempDir, "file2.thrift", file2Content})

	createTestFiles(t, tcArr)

	dynamicMap := &DynamicMap{
		innerMap: make(map[string]genericclient.Client),
	}

	testCases := []struct {
		name      string
		idlPath   string
		wantError bool
	}{
		{
			name:    "file1.thrift",
			idlPath: file1Path,
			// expect no error to be returned.
			wantError: false,
		},
		{
			name:    "file2.thrift",
			idlPath: file2Path,
			// expect errors to be returned.
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := dynamicMap.Add(tempDir)
			gotError := (err != nil)
			if tc.wantError != gotError {
				t.Errorf("Got error = %v, wantErr %v", gotError, tc.wantError)
				return
			} else if !tc.wantError {
				assert.True(t, dynamicMap.hasService(strings.ReplaceAll(tc.name, ".thrift", "")))
			}
		})
	}
}

// MockMap implements the IMap interface for testing. It mocks the DynamicMap struct.
type MockMap struct {
	mock.Mock
}

func (m *MockMap) Add(idlPath string, opts ...client.Option) error {
	args := m.Called(idlPath, opts)
	return args.Error(0)
}

func (m *MockMap) Delete(idlFileName string) {
	m.Called(idlFileName)
	return
}

func TestDynamicMap_AddAll(t *testing.T) {
	// Create a temp directory with temp files for testing
	tempDir := "./tempdir"
	tempDir, err := os.MkdirTemp("", "temp")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	subDir := tempDir + "/subdir"
	err = os.MkdirAll(subDir, os.ModePerm)

	if err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	tcArr := []TestCase{}
	tcArr = append(tcArr, TestCase{tempDir, "file1.txt", []byte(``)})
	tcArr = append(tcArr, TestCase{tempDir, "file2.psd", []byte(``)})
	tcArr = append(tcArr, TestCase{subDir, "file3.html", []byte(``)})
	tcArr = append(tcArr, TestCase{subDir, "file4", []byte(``)})
	tcArr = append(tcArr, TestCase{tempDir, "file5.thrift", []byte(``)})
	createTestFiles(t, tcArr)

	// Use a MockMap to assert expectations by AddAll(...)
	mockMap := &MockMap{
		Mock: mock.Mock{},
	}

	mockMap.On("Add", "file1.txt", tempDir, mock.Anything).Return(nil).Once()
	mockMap.On("Add", "file2.psd", tempDir, mock.Anything).Return(nil).Once()
	mockMap.On("Add", "file3.html", subDir, mock.Anything).Return(nil).Once()
	mockMap.On("Add", "file4", subDir, mock.Anything).Return(nil).Once()
	mockMap.On("Add", "file5.thrift", tempDir, mock.Anything).Return(nil).Once()

	AddAll(mockMap, tempDir)

	mockMap.AssertExpectations(t)
}
