package idl_mapping

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
)

// an interface created for testing purposes.
type IMap interface {
	Add(idlPath string, opts ...client.Option) error
	Delete(idlFileName string)
}

type DynamicMap struct {
	// define a map from the service name (given by the IDL file)
	// to the corresponding kitex generic client.
	innerMap map[string]genericclient.Client
}

func (m *DynamicMap) GetClient(svcName string) (cli genericclient.Client, ok bool) {
	cli, ok = m.innerMap[svcName]
	return cli, ok
}

func (m *DynamicMap) Add(idlPath string, opts ...client.Option) error {
	idlFileName := filepath.Base(idlPath)
	svcName := strings.ReplaceAll(idlFileName, ".thrift", "")

	// creating a new generic client.
	p, err := generic.NewThriftFileProvider(idlPath)
	if err != nil {
		log.Printf("creating new thrift file provider failed: %v", err)
		return err
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		log.Printf("creating new generic instance failed: %v", err)
		return err
	}

	cli, err := genericclient.NewClient(
		svcName,
		g,
		opts...,
	)
	if err != nil {
		log.Printf("creating new generic client failed: %v", err)
		return err
	}

	// adding the generic client to the map.
	if m.innerMap == nil {
		m.innerMap = make(map[string]genericclient.Client)
	}
	m.innerMap[svcName] = cli
	return nil
}

func (m *DynamicMap) Delete(svcName string) {
	delete(m.innerMap, svcName)
}

func AddAll(m IMap, idlPath string, opts ...client.Option) {
	c, err := os.ReadDir(idlPath)
	if err != nil {
		log.Fatalf("scanning idl file directory failed: %v", err)
	}

	for _, entry := range c {
		if entry.IsDir() {
			AddAll(m, idlPath+"/"+entry.Name(), opts...)
		} else {
			m.Add(idlPath+"/"+entry.Name(), opts...)
		}
	}
}
