package main

import (
	"github.com/shlande/ddns"
	"github.com/shlande/ddns/provider"
	"github.com/sirupsen/logrus"
)

type BindConfig struct {
	SecretKey string
	SecretId  string
	Provider  string

	Domain string
	Prefix string
}

func buildBinder(config BindConfig) *ddns.Binder {
	return ddns.NewBinder(
		config.Prefix,
		buildDomain(config.Provider, config.Domain, config.SecretId, config.SecretKey),
	)
}

func buildDomain(tp, domain, id, key string) ddns.Domain {
	var dm ddns.Domain
	var err error
	switch tp {
	case "alidns":
		dm, err = provider.NewAliDomain(domain, id, key)
	case "dnspod":
		dm, err = provider.NewDnsPodDomain(domain, id, key)
	default:
		logrus.Fatal("无效的dns类型,目前只支持alidns和dnspod")
	}
	if err != nil {
		logrus.Fatalf("无法创建dns:%v", err)
	}
	return dm
}
