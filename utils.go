package ddns

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
)

func GetIP() (string, error) {
	var status, ip string
	rsp := make(map[string]interface{})
	request, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(data, &rsp)
	if err != nil {
		return "", err
	}
	if s, ok := rsp["status"].(string); ok {
		status = s
	}
	if quey, ok := rsp["query"].(string); ok {
		ip = quey
	}
	if status != "success" {
		return "", errors.New("server error")
	}
	return ip, nil
}

func GetXdIp() (string, error) {
	resp, err := http.DefaultClient.Get("https://linux.xidian.edu.cn/ip")
	if err != nil {
		return "", err
	}
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(bt), "\n", ""), nil
}

func GetIPv6() (string, error) {
	s, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, a := range s {
		i := regexp.MustCompile(`(\w+:){7}\w+`).FindString(a.String())
		if strings.Count(i, ":") == 7 {
			return i, nil
		}
	}
	return "", errors.New("未找到有效ipv6地址")
}
