package ddns

import (
	"testing"
)

func TestNewDnsPod(t *testing.T) {
	client := NewDnsPod(
		"",
		"",
		DomainInfo{
			Type:       "xd",
			DomainName: "",
			Prefix:     "",
		},
	)
	err := client.Update()
	if err != nil {
		panic(err)
	}
}
