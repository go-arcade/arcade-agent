package main

import (
	"flag"

	"github.com/go-arcade/arcade-agent/internal/grpc"
	"github.com/go-arcade/arcade-agent/internal/grpc/conf"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "conf", "conf.d/config.toml", "conf file path, e.g. -conf ./conf.d")
}

func main() {
	flag.Parse()

	// 加载配置
	conf.NewConf(configFile)

	grpc.NewGrpcClient("localhost:9090")

}
