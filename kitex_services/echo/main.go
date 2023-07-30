package main

import (
	"log"
	"net"
	"sync"

	"github.com/Linda-ui/orbital_HeBao/kitex_services/echo/config"
	handler "github.com/Linda-ui/orbital_HeBao/kitex_services/echo/handler"
	echosvc "github.com/Linda-ui/orbital_HeBao/kitex_services/echo/kitex_gen/echo/echosvc"
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

			svr := echosvc.NewServer(
				new(handler.EchoImpl),
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
