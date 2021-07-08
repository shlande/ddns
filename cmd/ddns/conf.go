package main

import (
	"ddns"
	"errors"
)

func ParseAliDNS(args []string) (ddns.DNS, error) {
	if len(args) != 2 {
		panic(errors.New("错误的参数个数"))
	}
	keyId, secret := args[0], args[1]
	return ddns.NewAliDNS(keyId, secret, ddns.DomainInfo{
		Type:       tp,
		DomainName: domain,
		Prefix:     prefix,
	})
}

func ParseDnsPod(args []string) (ddns.DNS, error) {
	if len(args) != 2 {
		panic(errors.New("错误的参数个数"))
	}
	keyId, secret := args[0], args[1]
	return ddns.NewDnsPod(keyId, secret, ddns.DomainInfo{
		Type:       tp,
		DomainName: domain,
		Prefix:     prefix,
	}), nil
}
