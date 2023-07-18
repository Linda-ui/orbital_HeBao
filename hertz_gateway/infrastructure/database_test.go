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

	root, err := test.GetIDLRoot()
	if err != nil {
		t.Fatalf("failed to get IDL directory: %v", err)
	}

	file1_path := filepath.Join(*root, "file1.thrift")
	file2_path := filepath.Join(*root, "file2.thrift")

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
			args:      args{file1_path},
			svcName:   "file1",
			wantError: false,
		},
		{
			Name:      "invalid thrift IDL file",
			args:      args{file2_path},
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
