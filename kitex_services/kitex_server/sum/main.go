package main

import (
	"log"
	"net"
	"sync"

	"github.com/Linda-ui/orbital_HeBao/kitex_services/kitex_gen/sum/sumsvc"
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

	ports := []int{8880, 8881, 8882}

	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)

		go func(p int) {
			addr := &net.TCPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: p,
			}

			svr := sumsvc.NewServer(
				new(handler.SumImpl),
				server.WithRegistry(r),
				server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "sum"}),
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
