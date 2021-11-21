package provider

import (
	"fmt"
	"github.com/shlande/ddns"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"strconv"
)

var recordLine = common.StringPtr("默认")

type dnsPodClient struct {
	*dnspod.Client
	domains map[string]uint64
}

func NewDnsPodClient(secretId, secretKey string) (*dnsPodClient, error) {
	credential := common.NewCredential(secretId, secretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, _ := dnspod.NewClient(credential, "", cpf)

	request := dnspod.NewDescribeDomainListRequest()

	response, err := client.DescribeDomainList(request)
	if err != nil {
		return nil, err
	}
	// 获取文件信息
	var domains = make(map[string]uint64)
	for _, v := range response.Response.DomainList {
		domains[*v.Name] = *v.DomainId
	}
	return &dnsPodClient{
		Client:  client,
		domains: domains,
	}, nil
}

func NewDnsPodDomain(domain, secretId, secretKey string) (*dnsPodDomain, error) {
	client, err := NewDnsPodClient(secretId, secretKey)
	if err != nil {
		return nil, err
	}
	return newDnsPodDomain(domain, client)
}

func newDnsPodDomain(domain string, client *dnsPodClient) (*dnsPodDomain, error) {
	if id, has := client.domains[domain]; has {
		return &dnsPodDomain{
			dnsPodClient: client,
			domain:       domain,
			id:           strconv.Itoa(int(id)),
			id_:          id,
		}, nil
	}
	return nil, fmt.Errorf("用户账号下没有该域名: %v", domain)
}

type dnsPodDomain struct {
	*dnsPodClient
	// 类型
	domain string
	// 域名id
	id  string
	id_ uint64
}

func (d *dnsPodDomain) GetByPrefix(prefix string) ([]*ddns.Record, error) {
	request := dnspod.NewDescribeRecordListRequest()

	request.DomainId = common.Uint64Ptr(d.id_)
	request.Domain = common.StringPtr(d.domain)
	request.Subdomain = common.StringPtr(prefix)

	response, err := d.dnsPodClient.DescribeRecordList(request)
	if err != nil {
		return nil, err
	}
	var rcds = make([]*ddns.Record, 0, len(response.Response.RecordList))
	for _, v := range response.Response.RecordList {
		rcds = append(rcds, &ddns.Record{
			Type:       *v.Type,
			DomainName: d.domain,
			RecordID:   strconv.Itoa(int(*v.RecordId)),
			Value:      *v.Value,
			Prefix:     *v.Name,
		})
	}
	return rcds, nil
}

func (d *dnsPodDomain) CreateByRecords(rcds ...*ddns.Record) error {
	request := dnspod.NewCreateRecordRequest()

	request.Domain = common.StringPtr(d.domain)
	request.DomainId = common.Uint64Ptr(d.id_)
	for _, v := range rcds {
		request.Value = common.StringPtr(v.Value)
		request.RecordType = common.StringPtr(v.Type)
		request.SubDomain = common.StringPtr(v.Prefix)
		request.RecordLine = recordLine
		resp, err := d.Client.CreateRecord(request)
		if err != nil {
			return err
		}
		v.RecordID = strconv.Itoa(int(*resp.Response.RecordId))
	}
	return nil
}

func (d *dnsPodDomain) UpdateByRecords(rcds ...*ddns.Record) error {
	request := dnspod.NewModifyRecordRequest()

	request.Domain = common.StringPtr(d.domain)
	request.RecordLine = recordLine
	for _, v := range rcds {
		id, _ := strconv.ParseUint(v.RecordID, 10, 64)
		request.RecordId = common.Uint64Ptr(id)
		request.RecordType = common.StringPtr(v.Type)
		request.Value = common.StringPtr(v.Value)

		_, err := d.Client.ModifyRecord(request)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *dnsPodDomain) DeleteByRecords(rcds ...*ddns.Record) error {
	request := dnspod.NewDeleteRecordRequest()
	request.Domain = common.StringPtr(d.domain)

	for _, v := range rcds {
		id, _ := strconv.ParseUint(v.RecordID, 10, 64)
		request.RecordId = common.Uint64Ptr(id)
		_, err := d.Client.DeleteRecord(request)
		if err != nil {
			return err
		}
	}
	return nil
}
