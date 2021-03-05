package ddns

import (
	"fmt"
	"testing"
)

func createAliDNS() *AliDNS {
	dns, err := NewAliDNS("", "",
		DomainInfo{
			Type:       "",
			DomainName: "",
			Prefix:     "",
		},
	)
	if err != nil {
		panic(err)
	}
	return dns.(*AliDNS)
}

func TestAliDNS_describeRecords(t *testing.T) {
	c := createAliDNS()
	rcds, err := c.describeRecords()
	if err != nil {
		panic(err)
	}
	fmt.Println(rcds)
}

func TestAliDNS_Update(t *testing.T) {
	c := createAliDNS()
	err := c.Update()
	if err != nil {
		panic(err)
	}
}
