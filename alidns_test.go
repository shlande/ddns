package ddns

import (
	"fmt"
	"testing"
)

func createAliDNS() *AliDNS {
	dns,err := NewAliDNS("","",
		DomainInfo{
			DomainName: "",
			Prefix: "",
		},
	)
	if err != nil {
		panic(err)
	}
	return dns.(*AliDNS)
}

func TestAliDNS_describeRecords(t *testing.T) {
	c := createAliDNS()
	rcds,err := c.describeRecords()
	if err != nil {
		panic(err)
	}
	fmt.Println(rcds)
}

func TestAliDNS_addRecords(t *testing.T) {
	c := createAliDNS()
	c.addRecords()
}

func TestAliDNS_updateRecord(t *testing.T) {
	c := createAliDNS()
	err := c.updateRecords("19817912634469376")
	if err != nil {
		panic(err)
	}
}

func TestAliDNS_findRecordID(t *testing.T) {
	c := createAliDNS()
	rcds,err := c.describeRecords()
	if err != nil {
		panic(err)
	}
	if id := c.findRecordID(rcds); id!= "19817912634469376" {
		panic("错误的id："+ id)
	}
}

func TestAliDNS_Update(t *testing.T) {
	c := createAliDNS()
	err := c.Update()
	if err != nil {
		panic(err)
	}
}