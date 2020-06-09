package ddns

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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
