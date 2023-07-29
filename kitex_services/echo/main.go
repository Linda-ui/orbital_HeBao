package main

import (
	"log"
	"net"
	"path/filepath"
	"runtime"
	"sync"

	handler "github.com/Linda-ui/orbital_HeBao/kitex_services/echo/handler"
	echosvc "github.com/Linda-ui/orbital_HeBao/kitex_services/echo/kitex_gen/echo/echosvc"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
	"github.com/spf13/viper"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		klog.Fatal("failed to set cwd to echo service root directory")
	}

	viper.AddConfigPath(filepath.Dir(filename))
	viper.SetConfigName("server_config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		klog.Fatal(err)
	}

	ports := viper.GetStringSlice("ports")

	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)

		go func(p string) {
			addr, _ := net.ResolveTCPAddr("tcp", viper.GetString("host")+":"+p)

			svr := echosvc.NewServer(
				new(handler.EchoImpl),
				server.WithRegistry(r),
				server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: viper.GetString("serviceName")}),
				server.WithServiceAddr(addr),
			)

			err := svr.Run()
			if err != nil {
				log.Println(err)
			}
		}(port)
	}
	wg.Wait()
}
