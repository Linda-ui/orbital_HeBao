package idl_mapping

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	// Make a temporary directory
	tempDir, err := os.MkdirTemp("", "dir")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create and write the temporary IDL files for testing
	file0Path := filepath.Join(tempDir, "file0.thrift")
	file1Path := filepath.Join(tempDir, "file1.thrift")
	file2Path := filepath.Join(tempDir, "file2.thrift")
	file3Path := filepath.Join(tempDir, "file3.thrift")

	tcArr := []TestCase{}

	file0Content := []byte(``)
	tcArr = append(tcArr, TestCase{tempDir, "file0.thrift", file0Content})

	file1Content := []byte(`
		namespace go example.file1

		struct MyStruct {
			1: required string name,
			2: optional i32 age,
		}
	`)
	tcArr = append(tcArr, TestCase{tempDir, "file1.thrift", file1Content})

	file2Content := []byte(`
		namespace go example.file2

		service MyService {
			i32 add(1: i32 a, 2: i32 b),
		}
	`)
	tcArr = append(tcArr, TestCase{tempDir, "file2.thrift", file2Content})

	file3Content := []byte(`
		namespace go example.subdir.file3

		enum MyEnum {
			ONE,
			TWO,
			THREE,
		}

		service MyService {
			i32 add(1: i32 a, 2: i32 b),
		}
	`)
	tcArr = append(tcArr, TestCase{tempDir, "file3.thrift", file3Content})

	createTestFiles(t, tcArr)

	// Write testcases
	dynamicMap := &DynamicMap{
		innerMap: make(map[string]genericclient.Client),
	}

	testCases := []struct {
		name      string
		idlPath   string
		expectErr string
	}{
		{
			name:      "file0.thrift",
			idlPath:   file0Path,
			expectErr: fmt.Sprintf("parse ../../../../../../..%s/file0.thrift err: not document", tempDir),
		},
		{
			name:      "file1.thrift",
			idlPath:   file1Path,
			expectErr: "empty serverce from idls",
		},
		{
			name:      "file2.thrift",
			idlPath:   file2Path,
			expectErr: "",
		},
		{
			name:      "file3.thrift",
			idlPath:   file3Path,
			expectErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := dynamicMap.Add(tc.name, tempDir)
			switch tc.name {
			case "file0.thrift":
				assert.EqualError(t, err, tc.expectErr)
				assert.False(t, dynamicMap.hasService("file0"))
			case "file1.thrift":
				assert.EqualError(t, err, tc.expectErr)
				assert.False(t, dynamicMap.hasService("file1"))
			default:
				assert.NoError(t, err)
			}
		})
	}
	assert.True(t, dynamicMap.hasService("file2"))
	assert.True(t, dynamicMap.hasService("file3"))
}

type MockMap struct {
	mock.Mock
	*DynamicMap
}

func (m *MockMap) Add(idlFileName string, idlPath string, opts ...client.Option) error {
	args := m.Called(idlFileName, idlPath, opts)
	return args.Error(0)
}

func TestDynamicMap_AddAll(t *testing.T) {
	// file 1 and 2 is for testing files in the main directory
	// file 3 and 4 is for testing files in the subdirectory

	tempDir, err := os.MkdirTemp("", "dir")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	subDir := filepath.Join(tempDir, "subdir")
	err = os.MkdirAll(subDir, os.ModePerm)

	if err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	tcArr := []TestCase{}
	tcArr = append(tcArr, TestCase{tempDir, "file1.thrift", []byte(``)})
	tcArr = append(tcArr, TestCase{tempDir, "file2.thrift", []byte(``)})
	tcArr = append(tcArr, TestCase{subDir, "file3.thrift", []byte(``)})
	tcArr = append(tcArr, TestCase{subDir, "file4.thrift", []byte(``)})
	createTestFiles(t, tcArr)

	mockDynamicMap := &DynamicMap{
		innerMap: make(map[string]genericclient.Client),
	}

	mockMap := &MockMap{
		Mock:       mock.Mock{},
		DynamicMap: mockDynamicMap,
	}

	mockMap.On("Add", "file1.txt", tempDir).Return(nil).Once()
	mockMap.On("Add", "file2.psd", tempDir).Return(nil).Once()
	mockMap.On("Add", "file3.html", subDir).Return(nil).Once()
	mockMap.On("Add", "file4", subDir).Return(nil).Once()

	mockMap.AddAll(tempDir)

	mockMap.AssertExpectations(t)
}
