package main

import (
	"log"
	"net"
	"sync"

	"github.com/Linda-ui/orbital_HeBao/kitex_services/sum/config"
	handler "github.com/Linda-ui/orbital_HeBao/kitex_services/sum/handler"
	"github.com/Linda-ui/orbital_HeBao/kitex_services/sum/kitex_gen/sum/sumsvc"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
)

func main() {
	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		klog.Fatal(err)
	}

	addrs := config.ServiceAddrs

	var wg sync.WaitGroup

	for _, addr := range addrs {
		wg.Add(1)

		go func(addr string) {
			netAddr, _ := net.ResolveTCPAddr("tcp", addr)

			svr := sumsvc.NewServer(
				new(handler.SumImpl),
				server.WithRegistry(r),
				server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.ServiceName}),
				server.WithServiceAddr(netAddr),
			)

			err := svr.Run()
			if err != nil {
				log.Println(err)
			}
		}(addr)
	}
	wg.Wait()
}
