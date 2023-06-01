package main

import (
	echosvc "Orbital_Hebao/kitex_servers/kitex_gen/echo/echosvc"
	"Orbital_Hebao/kitex_servers/kitex_gen/sum/sumsvc"
	handler "Orbital_Hebao/kitex_servers/kitex_handler"
	"log"

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

	svr1 := echosvc.NewServer(
		new(handler.EchoImpl),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "echo"}),
	)

	err1 := svr1.Run()
	if err1 != nil {
		log.Println(err1.Error())
	}

	svr2 := sumsvc.NewServer(
		new(handler.SumImpl),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "sum"}),
	)

	err2 := svr2.Run()
	if err2 != nil {
		log.Println(err2.Error())
	}
}
