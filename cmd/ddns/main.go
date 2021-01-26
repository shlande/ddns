package main

import (
	"ddns"
	"flag"
	"log"
	"time"
)

var provider, domain, prefix string

var ttl, maxRetry int

func main() {
	var (
		dns ddns.DNS
		err error
	)
	flag.StringVar(&provider, "provider", "", "ddns服务提供商，目前只支持alidns")
	flag.StringVar(&domain, "domain", "", "需要解析的域名")
	flag.StringVar(&prefix, "prefix", "", "域名的前缀")
	flag.IntVar(&ttl, "ttl", 30, "查询间隔")
	flag.IntVar(&maxRetry, "retry", 60, "最大连续出错次数")

	flag.Parse()

	switch provider {
	case "alidns":
		dns, err = ParseAliDNS(flag.Args())
	default:
		flag.PrintDefaults()
		return
	}

	if err != nil {
		panic(err)
	}
	Run(dns)
}

func Run(dns ddns.DNS) {
	// 立刻测试，如果出错则panic
	err := dns.Update()
	if err != nil {
		panic(err)
	}
	var curErrTime int
	for {
		time.Sleep(time.Duration(ttl) * time.Second)
		err := dns.Update()
		if err != nil {
			log.Println(err)
			curErrTime++
			if curErrTime > maxRetry {
				panic("连续出错次数超过允许的范围")
			}
			continue
		}
		curErrTime = 0
	}
}
