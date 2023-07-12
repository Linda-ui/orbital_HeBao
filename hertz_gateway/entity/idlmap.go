package entity

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
)

type MapManager interface {
	GetClient(serviceName string) (genericclient.Client, bool)
	AddService(idlPath string, opts ...client.Option) error
	DeleteService(serviceName string)
	AddAllServices(idlRootPath string, opts ...client.Option)
	DynamicUpdate(idlRootPath string, opts ...client.Option)
}
