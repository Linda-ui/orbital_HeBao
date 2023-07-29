package infrastructure

import (
	"log"

	utils "github.com/Linda-ui/orbital_HeBao/hertz_gateway/utils"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
)

// type database implements the idlmap.Repository interface.
type database map[string]genericclient.Client

func NewDatabase() database {
	return make(map[string]genericclient.Client)
}

func (db database) GetClient(svcName string) (genericclient.Client, bool) {
	cli, ok := db[svcName]
	return cli, ok
}

func (db database) AddService(idlPath string, opts ...client.Option) error {
	// extracting the IDL file name without extension as the service name
	svcName := utils.ExtractServiceName(idlPath)

	// creating a new generic client.
	p, err := generic.NewThriftFileProvider(idlPath)
	if err != nil {
		log.Printf("creating new thrift file provider for %v service failed: %v", svcName, err)
		return err
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		log.Printf("creating new generic instance for %v service failed: %v", svcName, err)
		return err
	}

	cli, err := genericclient.NewClient(
		svcName,
		g,
		opts...,
	)
	if err != nil {
		log.Printf("creating new generic client for %v service failed: %v", svcName, err)
		return err
	}

	// adding the generic client to the map.
	if db == nil {
		db = make(map[string]genericclient.Client)
	}
	db[svcName] = cli
	return nil
}

func (db database) DeleteService(svcName string) {
	delete(db, svcName)
}
