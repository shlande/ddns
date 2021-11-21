package provider

import (
	"fmt"
	"testing"
)

func getTestCloudflare() *cloudflare {
	cf, err := NewCloudFlare(Token)
	if err != nil {
		panic(err)
	}
	return cf
}

func TestCloudflare_listZone(t *testing.T) {
	resp, err := getTestCloudflare().listZone()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func TestCloudFlare_listRecords(t *testing.T) {
	res, err := getTestCloudflare().listRecords("shining.moe", "api")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestCloufFlare_CreateRecords(t *testing.T) {
	rcd := &cloudflareRecord{
		Type:    "A",
		Name:    "test.shining.moe",
		Content: "10.0.0.1",
	}
	err := getTestCloudflare().createRecords(rcd)
	if err != nil {
		panic(err)
	}
	fmt.Println(rcd)
}

func TestCloudflare_updateRecord(t *testing.T) {
	rcd := &cloudflareRecord{
		Id:      "043cbd794120da9c380446c60fb55d25",
		Type:    "A",
		Name:    "test.shining.moe",
		Content: "10.0.0.2",
	}
	err := getTestCloudflare().updateRecords(rcd)
	if err != nil {
		panic(err)
	}
	fmt.Println(rcd)
}

func TestCloudflare_deleteRecord(t *testing.T) {
	rcd := &cloudflareRecord{
		Id:      "043cbd794120da9c380446c60fb55d25",
		Type:    "A",
		Name:    "test.shining.moe",
		Content: "10.0.0.2",
	}
	err := getTestCloudflare().deleteRecords(rcd)
	if err != nil {
		panic(err)
	}
	fmt.Println(rcd)
}
