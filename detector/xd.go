package detector

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Xd struct{}

func (x Xd) IP() ([]string, error) {
	resp, err := http.DefaultClient.Get("https://linux.xidian.edu.cn/ip")
	if err != nil {
		return nil, err
	}
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return []string{strings.ReplaceAll(string(bt), "\n", "")}, nil
}
