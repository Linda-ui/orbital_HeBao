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

// MockClient implements the genericclient.Client interface for testing
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

func (m *DynamicMap) hasService(serviceName string) bool {
	_, ok := m.innerMap[serviceName]
	return ok
}

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

type TestCase struct {
	path    string
	name    string
	content []byte
}

func TestDynamicMap_Add(t *testing.T) {
	// Create temp directory with temp files for testing
	tempDir, err := os.MkdirTemp("", "dir")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	file1Path := filepath.Join(tempDir, "file1.thrift")
	file2Path := filepath.Join(tempDir, "file2.thrift")
	tcArr := []TestCase{}

	file1Content := []byte(`
		namespace go example.file1

		service MyService {
			i32 add(1: i32 a, 2: i32 b),
		}
	`)
	tcArr = append(tcArr, TestCase{tempDir, "file1.thrift", file1Content})

	file2Content := []byte(``)
	tcArr = append(tcArr, TestCase{tempDir, "file2.thrift", file2Content})

	createTestFiles(t, tcArr)

	dynamicMap := &DynamicMap{
		innerMap: make(map[string]genericclient.Client),
	}

	testCases := []struct {
		name    string
		idlPath string
		ok      bool
	}{
		{
			name:    "file1.thrift",
			idlPath: file1Path,
			ok:      true,
		},
		{
			name:    "file2.thrift",
			idlPath: file2Path,
			ok:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok := dynamicMap.Add(tc.name, tempDir)
			if ok != tc.ok {
				t.Errorf("Got error = %v, wantErr %v", ok, tc.ok)
				return
			} else if tc.ok {
				assert.True(t, dynamicMap.hasService(strings.ReplaceAll(tc.name, ".thrift", "")))
			}
		})
	}
}

type MockMap struct {
	mock.Mock
}

func (m *MockMap) Add(idlFileName string, idlPath string, opts ...client.Option) bool {
	m.Called(idlFileName, idlPath, opts)
	return true
}

func TestDynamicMap_AddAll(t *testing.T) {
	// Create temp directory with temp files for testing
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
