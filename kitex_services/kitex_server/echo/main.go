package main

import (
	"log"
	"net"
	"sync"

	echosvc "github.com/Linda-ui/orbital_HeBao/kitex_services/kitex_gen/echo/echosvc"
	handler "github.com/Linda-ui/orbital_HeBao/kitex_services/kitex_handler"

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

	ports := []int{8870, 8871, 8872}

	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)

		go func(p int) {
			addr := &net.TCPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: p,
			}

			svr := echosvc.NewServer(
				new(handler.EchoImpl),
				server.WithRegistry(r),
				server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "echo"}),
				server.WithServiceAddr(addr),
			)

			err := svr.Run()
			if err != nil {
				log.Println(err)
			}
		}(port)
	}
	wg.Wait()

	select {}

}
