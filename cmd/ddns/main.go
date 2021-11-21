package main

import (
	"context"
	"flag"
	"github.com/shlande/ddns"
)

type Config struct {
	BindConfig
	// 制定类型，ip或者ipv4
	Type   string
	Detect string
	TTL    int
}

var config Config

func main() {
	flag.StringVar(&config.Provider, "provider", "", "ddns服务提供商，目前支持alidns和dnspod")
	flag.StringVar(&config.Domain, "domain", "", "需要解析的域名")
	flag.StringVar(&config.Prefix, "prefix", "", "域名的前缀")
	flag.IntVar(&config.TTL, "ttl", 30, "查询间隔")
	flag.StringVar(&config.Detect, "detect", "ip", "绑定ip获取器，支持public，xd, device")
	flag.StringVar(&config.Type, "type", "ip", "网络类型，支持ip，ipv6")

	flag.Parse()

	binder := buildBinder(config.BindConfig)

	switch config.Type {
	case "ip":
		detector := buildDetector(config.Detect)
		ddns.RunIpWorker(context.Background(), detector, binder, int64(config.TTL))
	case "ipv6":
		detector := buildDetectorV6(config.Detect)
		ddns.RunIpv6Worker(context.Background(), detector, binder, int64(config.TTL))
	}
}
