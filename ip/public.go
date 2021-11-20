package ip

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// PublicIpGetter 获取外网ip
type PublicIpGetter struct{}

func (p PublicIpGetter) IP() ([]string, error) {
	var status, ip string
	rsp := make(map[string]interface{})
	request, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &rsp)
	if err != nil {
		return nil, err
	}
	if s, ok := rsp["status"].(string); ok {
		status = s
	}
	if quey, ok := rsp["query"].(string); ok {
		ip = quey
	}
	if status != "success" {
		return nil, errors.New("server error")
	}
	return []string{ip}, nil
}
