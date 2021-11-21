package provider

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/shlande/ddns"
)

type AliAccount struct {
	Name      string
	SecretKey string
	SecretId  string
}

func NewAliDomain(domain, id, key string) (*AliDomain, error) {
	client, err := alidns.NewClientWithAccessKey("", id, key)
	if err != nil {
		return nil, err
	}
	// TODO：查询id
	return &AliDomain{
		client: client,
		id:     "",
		domain: domain,
	}, nil
}

// AliDomain 管理一个根域名信息
type AliDomain struct {
	client *alidns.Client
	// 域名id
	id string
	// 域名
	domain string
}

// GetByPrefix 通过前缀获取所有记录值
func (a *AliDomain) GetByPrefix(prefix string) ([]*ddns.Record, error) {
	req := alidns.CreateDescribeDomainRecordsRequest()
	req.Scheme = "https"
	req.DomainName = a.domain
	// TODO: 这里可能会导致问题
	req.SearchMode = "EXACT"
	req.KeyWord = prefix
	// TODO: 尝试支持500条以上的自动分页查询
	req.PageSize = "500"

	resp, err := a.client.DescribeDomainRecords(req)
	if err != nil {
		return nil, err
	}
	// 然后刷新结果
	var rcds = make([]*ddns.Record, 0, len(resp.DomainRecords.Record))
	for _, rcd := range resp.DomainRecords.Record {
		if rcd.RR == prefix {
			rcds = append(rcds, &ddns.Record{
				Type:       rcd.Type,
				DomainName: rcd.DomainName,
				RecordID:   rcd.RecordId,
				Value:      rcd.Value,
				Prefix:     rcd.RR,
			})
		}
	}
	return rcds, nil
}

// CreateByRecords 创建记录
func (a *AliDomain) CreateByRecords(rcds ...*ddns.Record) error {
	req := alidns.CreateAddDomainRecordRequest()
	req.Scheme = "https"

	for _, rcd := range rcds {
		req.DomainName = a.domain
		req.RR = rcd.Prefix
		req.Value = rcd.Value
		req.Type = rcd.Type

		resp, err := a.client.AddDomainRecord(req)
		if err != nil {
			return err
		}
		rcd.RecordID = resp.RecordId
	}
	return nil
}

// UpdateByRecords 通过ddns.Record更新ip
func (a *AliDomain) UpdateByRecords(rcd ...*ddns.Record) error {
	req := alidns.CreateUpdateDomainRecordRequest()
	req.Scheme = "https"

	for _, v := range rcd {
		req.RecordId = v.RecordID
		req.RR = v.Prefix
		req.Value = v.Value
		req.Type = v.Type

		_, err := a.client.UpdateDomainRecord(req)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AliDomain) DeleteByRecords(rcd ...*ddns.Record) error {
	req := alidns.CreateDeleteDomainRecordRequest()
	req.Scheme = "https"

	for _, v := range rcd {
		req.RecordId = v.RecordID
		_, err := a.client.DeleteDomainRecord(req)
		if err != nil {
			return err
		}
	}
	return nil
}
