package ddns

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

// 保存记录信息
type records struct {
	DomainName string
	RecordID   string
	Value      string
	Prefix     string
}

type AliDNS struct {
	c  *alidns.Client
	dn string
	p  string
}

func (a *AliDNS) Update() error {
	rcds, err := a.describeRecords()
	if err != nil {
		return err
	}
	if id := a.findRecordID(rcds); id != "" {
		return a.updateRecords(id)
	}
	return a.addRecords()
}

func (a *AliDNS) findRecordID(r []records) string {
	for _, rcd := range r {
		if rcd.Prefix == a.p {
			return rcd.RecordID
		}
	}
	return ""
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

func (a *AliDNS) addRecords() error {
	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = a.dn
	request.RR = a.p
	request.Type = "A"
	request.Scheme = "https"
	ip, err := GetIP()
	if err != nil {
		return err
	}
	request.Value = ip
	_, err = a.c.AddDomainRecord(request)
	if err != nil {
		return err
	}
	return nil
}

func (a *AliDNS) updateRecords(id string) error {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = id
	request.RR = a.p
	request.Type = "A"
	ip, err := GetIP()
	if err != nil {
		return err
	}
	request.Value = ip
	_, err = a.c.UpdateDomainRecord(request)
	return err
}

func NewAliDNS(AccessKeyID string, AccessSecret string, info DomainInfo) (DNS, error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", AccessKeyID, AccessSecret)
	if err != nil {
		return nil, err
	}
	return &AliDNS{c: client, dn: info.DomainName, p: info.Prefix}, nil
}
