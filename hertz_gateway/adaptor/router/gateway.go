package router

import (
	handler "github.com/Linda-ui/orbital_HeBao/hertz_gateway/adaptor/handler"
	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/biz/idlmap"
	"github.com/Linda-ui/orbital_HeBao/hertz_gateway/infrastructure"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/loadbalance/lbcache"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

// registers the router of gateway
func RegisterGateway(r *server.Hertz) {
	nacosResolver, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		hlog.Fatalf("err:%v", err)
	}

	lb := loadbalance.NewWeightedBalancer()

	idlRootPath := "./idl"

	idlServiceMap := infrastructure.NewDatabase()
	dynamicMapManager := idlmap.NewManager(idlServiceMap)
	gateway := handler.NewGateway(dynamicMapManager)

	dynamicMapManager.AddAllServices(
		idlRootPath,
		client.WithResolver(nacosResolver),
		client.WithLoadBalancer(lb, &lbcache.Options{Cacheable: true}),
	)

	go dynamicMapManager.DynamicUpdate(
		idlRootPath,
		client.WithResolver(nacosResolver),
		client.WithLoadBalancer(lb, &lbcache.Options{Cacheable: true}),
	)

	group := r.Group("/gateway")
	group.POST("/:svc/:method", gateway.Handler)
}
