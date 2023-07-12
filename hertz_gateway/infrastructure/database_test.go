package infrastructure

import (
	"path/filepath"
	"testing"

	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/test"
	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestDatabase_GetClient(t *testing.T) {
	db := NewDatabase()
	mockClient := mock.NewClient()

	db["serviceA"] = mockClient

	client, ok := db.GetClient("serviceA")
	assert.True(t, ok)
	assert.Equal(t, mockClient, client)

	client, ok = db.GetClient("serviceB")
	assert.False(t, ok)
	assert.Nil(t, client)
}

func TestDatabase_AddService(t *testing.T) {
	// Create a temp directory with temp files for testing.
	tempDir := t.TempDir()

	// create test file 1.
	// It is a valid thrift file with a service. Client is expected to be created.
	file1 := test.TempFile{
		Name: "file1.thrift",
		Path: filepath.Join(tempDir, "file1.thrift"),
		Content: []byte(`
		namespace go example.file1

		service MyService {
			i64 add(1: i64 a, 2: i64 b)
		}
	`),
	}

	// create test file 2.
	// It is an invalid thrift file without a service. Error is expected to be returned.
	file2 := test.TempFile{
		Name:    "file2.thrift",
		Path:    filepath.Join(tempDir, "file2.thrift"),
		Content: []byte(``),
	}

	// create a list of test files.
	files := []test.TempFile{file1, file2}
	test.CreateTestFiles(t, files)

	db := NewDatabase()

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
