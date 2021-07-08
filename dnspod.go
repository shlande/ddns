package ddns

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"log"
)

func NewDnsPod(secretId, secretKey string, info DomainInfo) *DnsPod {
	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, err := dnspod.NewClient(credential, "", cpf)
	if err != nil {
		panic(err)
	}
	return &DnsPod{
		Client:     client,
		DomainInfo: info,
	}
}

type DnsPod struct {
	*dnspod.Client
	// 类型
	DomainInfo
}

func (d *DnsPod) Update() error {
	request := dnspod.NewModifyDynamicDNSRequest()

	ip, err := getIp(d.Type)
	if err != nil {
		panic(err)
	}

	// 获取域名列表
	recordId, value := d.getRecord(d.Prefix, d.DomainName)
	if recordId == 0 {
		recordId = d.createRecord(ip)
		return nil
	}
	if value == ip {
		return nil
	}

	request.Domain = common.StringPtr("colaha.tech")
	request.Value = common.StringPtr(ip)
	request.SubDomain = common.StringPtr(d.Prefix)
	request.RecordLine = common.StringPtr("默认")
	request.RecordId = common.Uint64Ptr(850058500)
	request.RecordId = common.Uint64Ptr(recordId)

	_, err = d.ModifyDynamicDNS(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	return err
}

func (d *DnsPod) getRecord(prefix, domain string) (recordId uint64, value string) {
	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = common.StringPtr(domain)

	response, err := d.DescribeRecordList(request)
	if err != nil {
		log.Println("无法获取记录信息:", err)
		return
	}
	for _, v := range response.Response.RecordList {
		if *v.Name == prefix {
			recordId = *v.RecordId
			value = *v.Value
		}
	}
	return
}

func (d *DnsPod) createRecord(ip string) uint64 {
	// 没有找到，尝试创建一个
	request := dnspod.NewCreateRecordRequest()
	request.Domain = common.StringPtr(d.DomainInfo.DomainName)
	switch d.Type {
	case "xd":
		fallthrough
	case "ipv4":
		request.RecordType = common.StringPtr("A")
	case "ipv6":
		request.RecordType = common.StringPtr("AAAA")
	default:
		panic("unknow type")
	}
	request.RecordLine = common.StringPtr("默认")
	request.Value = common.StringPtr(ip)
	request.SubDomain = common.StringPtr(d.Prefix)

	response, err := d.CreateRecord(request)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return 0
	}
	if err != nil {
		panic(err)
	}
	return *response.Response.RecordId
}
