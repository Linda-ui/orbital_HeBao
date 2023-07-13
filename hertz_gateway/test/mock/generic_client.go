package mock

import (
	"context"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/stretchr/testify/mock"
)

// fakeClient implements the genericclient.Client interface for testing
type fakeClient struct {
	mock.Mock
}

func NewClient() *fakeClient {
	return &fakeClient{}
}

func (c *fakeClient) GenericCall(ctx context.Context, method string, request interface{}, callOptions ...callopt.Option) (response interface{}, err error) {
	args := c.Called(ctx, method, request, callOptions)
	return args.Get(0), args.Error(1)
}

func (c *fakeClient) Close() error {
	args := c.Called()
	return args.Error(0)
}
