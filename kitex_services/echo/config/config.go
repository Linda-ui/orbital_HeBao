package config

import (
	"path/filepath"
	"runtime"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/spf13/viper"
)

var (
	ServiceName  string
	ServicePorts []string
	ServiceAddrs []string
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		klog.Fatal("failed to set cwd to echo service root directory")
	}

	viper.AddConfigPath(filepath.Dir(filename))
	viper.SetConfigName("server_config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	ServiceName = viper.GetString("serviceName")
	ServicePorts = viper.GetStringSlice("ports")
	ServiceAddrs = make([]string, len(ServicePorts))
	host := viper.GetString("host")
	for i, port := range ServicePorts {
		ServiceAddrs[i] = host + ":" + port
	}
}
