package main

import (
	"log"
	"net"

	echosvc "github.com/Linda-ui/orbital_HeBao/kitex_services/kitex_gen/echo/echosvc"
	handler "github.com/Linda-ui/orbital_HeBao/kitex_services/kitex_handler"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
)

func main() {
	r1, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		klog.Fatal(err)
	}

	svr1 := echosvc.NewServer(
		new(handler.EchoImpl),
		server.WithRegistry(r1),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "echo"}),
		server.WithServiceAddr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8888}),
	)

	err1 := svr1.Run()
	if err1 != nil {
		log.Println(err1.Error())
	}
}
