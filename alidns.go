package ddns

import (
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

// 保存记录信息
type records struct {
	Type       string
	DomainName string
	RecordID   string
	Value      string
	Prefix     string
}

type AliDNS struct {
	c  *alidns.Client
	dn string
	p  string
	tp string
}

func (a *AliDNS) Update() error {
	rcds, err := a.describeRecords()
	if err != nil {
		return err
	}
	ip, err := a.getIp()
	if err != nil {
		return err
	}
	if rcd := a.findRecordID(rcds); rcd.RecordID != "" {
		if rcd.Value == ip {
			return nil
		}
		return a.updateRecords(rcd.RecordID, ip)
	}
	return a.addRecords(ip)
}

func (a *AliDNS) getIp() (string, error) {
	var ip string
	var err error
	switch a.tp {
	case "ipv4":
		ip, err = GetIP()
	case "ipv6":
		ip, err = GetIPv6()
	case "xd":
		ip, err = GetXdIp()
	}
	return ip, err
}

func (a *AliDNS) findRecordID(r []records) records {
	for _, rcd := range r {
		if rcd.Prefix == a.p {
			return rcd
		}
	}
	return records{}
}

func (a *AliDNS) describeRecords() ([]records, error) {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.DomainName = a.dn
	request.Scheme = "https"
	response, err := a.c.DescribeDomainRecords(request)
	if err != nil {
		return nil, err
	}
	rcds := make([]records, 0, len(response.DomainRecords.Record))
	for _, rcd := range response.DomainRecords.Record {
		rcds = append(rcds, records{
			DomainName: rcd.DomainName,
			RecordID:   rcd.RecordId,
			Value:      rcd.Value,
			Prefix:     rcd.RR,
		})
	}
	return rcds, nil
}

func (a *AliDNS) addRecords(value string) error {
	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = a.dn
	request.RR = a.p
	if a.tp == "ipv4" || a.tp == "xd" {
		request.Type = "A"
	}
	if a.tp == "ipv6" {
		request.Type = "AAAA"
	}
	request.Scheme = "https"
	request.Value = value
	_, err := a.c.AddDomainRecord(request)
	if err != nil {
		return err
	}
	return nil
}

func (a *AliDNS) updateRecords(id string, value string) error {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = id
	request.RR = a.p
	request.Type = a.tp
	request.Value = value
	_, err := a.c.UpdateDomainRecord(request)
	return err
}

func NewAliDNS(AccessKeyID string, AccessSecret string, info DomainInfo) (DNS, error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", AccessKeyID, AccessSecret)
	if err != nil {
		return nil, err
	}
	if info.Type != "xd" && info.Type != "ipv4" && info.Type != "ipv6" {
		return nil, errors.New("无效的域名记录类型：" + info.Type)
	}
	return &AliDNS{c: client, dn: info.DomainName, p: info.Prefix, tp: info.Type}, nil
}
