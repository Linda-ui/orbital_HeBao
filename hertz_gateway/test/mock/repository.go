package mock

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/stretchr/testify/mock"
)

// fakeRepository implements the idlmap.Repository interface for testing
type fakeRepository struct {
	mock.Mock
}

func NewRepository() *fakeRepository {
	return &fakeRepository{}
}

func (repo *fakeRepository) GetClient(svcName string) (cli genericclient.Client, ok bool) {
	args := repo.Called(svcName)
	return args.Get(0).(genericclient.Client), args.Bool(1)
}

func (repo *fakeRepository) AddService(idlPath string, opts ...client.Option) error {
	args := repo.Called(idlPath, opts)
	return args.Error(0)
}

func (repo *fakeRepository) DeleteService(svcName string) {
	repo.Called(svcName)
}
