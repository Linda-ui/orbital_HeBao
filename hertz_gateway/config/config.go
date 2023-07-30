package config

import (
	"path/filepath"
	"runtime"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/spf13/viper"
)

var (
	ServiceAddr string
	ServiceName string
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		hlog.Fatal("failed to set cwd to hertz gateway root directory")
	}

	viper.AddConfigPath(filepath.Dir(filename))
	viper.SetConfigName("gateway_config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	ServiceAddr = viper.GetString("host") + ":" + viper.GetString("port")
	ServiceName = viper.GetString("serviceName")
}
