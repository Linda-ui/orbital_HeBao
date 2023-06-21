package idl_mapping

import (
	"log"
	"os"
	"strings"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
)

type DynamicMap struct {
	// define a map from the service name (given by the IDL file)
	// to the corresponding kitex generic client
	innerMap map[string]genericclient.Client
}

func (m *DynamicMap) GetClient(svcName string) (cli genericclient.Client, ok bool) {
	cli, ok = m.innerMap[svcName]
	return cli, ok
}

func (m *DynamicMap) AddAll(idlPath string, opts ...client.Option) {
	c, err := os.ReadDir(idlPath)
	if err != nil {
		log.Fatalf("scanning idl file directory failed: %v", err)
	}

	for _, entry := range c {
		if entry.IsDir() {
			m.AddAll(idlPath+entry.Name()+"/", opts...)
		} else {
			m.Add(entry.Name(), idlPath, opts...)
		}
	}
}

func (m *DynamicMap) Add(idlFileName string, idlPath string, opts ...client.Option) {
	svcName := strings.ReplaceAll(idlFileName, ".thrift", "")

	// creating a new generic client
	p, err := generic.NewThriftFileProvider(idlFileName, idlPath)
	if err != nil {
		log.Fatalf("creating new thrift file provider failed: %v", err)
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		log.Fatalf("creating new generic instance failed: %v", err)
	}

	cli, err := genericclient.NewClient(
		svcName,
		g,
		opts...,
	)
	if err != nil {
		log.Fatalf("creating new generic client failed: %v", err)
	}

	// adding the generic client to the map
	if m.innerMap == nil {
		m.innerMap = make(map[string]genericclient.Client)
	}
	m.innerMap[svcName] = cli
}
