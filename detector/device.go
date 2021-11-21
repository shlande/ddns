package detector

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"regexp"
	"strings"
)

var (
	ErrInterface = errors.New("无法获取网卡信息")
)

func NewDeviceGetter(match string) (*Device, error) {
	rex, err := regexp.Compile(match)
	if err != nil {
		return nil, fmt.Errorf("无法解析正则：%v", err)
	}
	// 首先尝试获取一次
	_, find, err := getIpFromDevice(rex, false, false)
	if err != nil {
		return nil, fmt.Errorf("无法获取网卡信息：%v", err)
	}
	if !find {
		logrus.Fatalf("没有找到匹配的网卡：%v", err)
	}
	return &Device{rex: rex}, nil
}

func getIpFromDevice(rex *regexp.Regexp, ipv4 bool, ipv6 bool) ([]*net.IPNet, bool, error) {
	ifac, err := net.Interfaces()
	if err != nil {
		return nil, false, fmt.Errorf("%w:%v", ErrInterface, err)
	}
	var ips = make([]*net.IPNet, 0, len(ifac))
	// 尝试正则匹配
	for _, v := range ifac {
		if !rex.MatchString(v.Name) {
			continue
		}
		addr, err := v.Addrs()
		if err != nil {
			return nil, true, fmt.Errorf("%w:%v", ErrInterface, err)
		}
		// 保存ip信息
		for _, v := range addr {
			ip := v.(*net.IPNet)
			// 如果ipv6，ipv4中
			if (ipv4 && ip.IP.To4() != nil) || ipv6 && ip.IP.To4() == nil {
				ips = append(ips, ip)
			}
		}
	}
	return ips, len(ifac) != 0, nil
}

type Device struct {
	rex *regexp.Regexp
}

func (d Device) IPv6() ([]string, error) {
	res, _, err := getIpFromDevice(d.rex, false, true)
	if err != nil {
		return nil, err
	}
	ips := make([]string, 0, len(res))
	for _, v := range res {
		ips = append(ips, unmask(v))
	}
	return ips, nil
}

func (d Device) IP() ([]string, error) {
	res, _, err := getIpFromDevice(d.rex, true, false)
	if err != nil {
		return nil, err
	}
	ips := make([]string, 0, len(res))
	for _, v := range res {
		ips = append(ips, unmask(v))
	}
	return ips, nil
}

func unmask(ipNet *net.IPNet) string {
	return strings.Split(ipNet.String(), "/")[0]
}
