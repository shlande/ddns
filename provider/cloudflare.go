package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shlande/ddns"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type zoneListResponse struct {
	Success bool
	Result  []struct {
		Id     string
		Name   string
		Status string
	}
}

type cloudflareRecord struct {
	Id       string
	Type     string
	Name     string
	Content  string
	TTL      int
	Priority int
	Proxied  bool
}

type listRecordResponse struct {
	Result []*cloudflareRecord
}

type recordResponse struct {
	Result cloudflareRecord
}

func NewCloudFlare(token string) (*cloudflare, error) {
	cf := &cloudflare{
		token:   token,
		domains: map[string]string{},
	}
	return cf, cf.cache()
}

type cloudflare struct {
	token string
	// key: name  value: id
	domains map[string]string
}

func (c *cloudflare) request(method, endpoint string, body interface{}, response interface{}) error {
	var (
		data []byte
		err  error
		dr   io.Reader
	)

	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return err
		}
		dr = bytes.NewReader(data)
	}

	request, err := http.NewRequest(method, fmt.Sprintf("https://api.cloudflare.com/client/v4/%v", endpoint), dr)
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if response != nil {
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, response)
	}
	return nil
}

func (c *cloudflare) listZone() (*zoneListResponse, error) {
	var resp = &zoneListResponse{}
	err := c.request("GET", "zones", nil, resp)
	if err != nil {
		return nil, fmt.Errorf("获取zones信息失败: %w", err)
	}
	return resp, nil
}

func (c *cloudflare) listRecords(zoneName, prefix string) (*listRecordResponse, error) {
	var name = zoneName
	if len(prefix) != 0 {
		name = prefix + "." + zoneName
	}
	if zoneId, has := c.domains[zoneName]; has {
		var resp = &listRecordResponse{}
		err := c.request("GET", fmt.Sprintf("zones/%v/dns_records?name=%v", zoneId, name), nil, resp)
		if err != nil {
			return nil, err
		}
		return resp, err
	}
	return nil, errors.New("无法找到域名信息")
}

func (c *cloudflare) createRecords(request *cloudflareRecord) error {
	aid, err := c.getZoneId(request.Name)
	if err != nil {
		return err
	}

	var resp = &recordResponse{}
	err = c.request(http.MethodPost, fmt.Sprintf("zones/%v/dns_records", aid), request, resp)
	if err != nil {
		return err
	}
	*request = resp.Result
	return nil
}

func (c *cloudflare) getZoneId(name string) (string, error) {
	raw := strings.Split(name, ".")
	var domain, aid string
	for i, _ := range raw {
		domain = strings.Join(raw[i:], ".")
		if areaId, has := c.domains[domain]; has {
			aid = areaId
			break
		}
	}
	if len(aid) == 0 {
		return "", errors.New("没有找到域名信息")
	}
	return aid, nil
}

func (c *cloudflare) updateRecords(record *cloudflareRecord) error {
	aid, err := c.getZoneId(record.Name)
	if err != nil {
		return err
	}
	var res = &recordResponse{}
	err = c.request("PUT", fmt.Sprintf("zones/%v/dns_records/%v", aid, record.Id), record, res)
	if err != nil {
		return err
	}
	*record = res.Result
	return nil
}

func (c *cloudflare) deleteRecords(record *cloudflareRecord) error {
	aid, err := c.getZoneId(record.Name)
	if err != nil {
		return err
	}
	return c.request("DELETE", fmt.Sprintf("zones/%v/dns_records/%v", aid, record.Id), nil, nil)
}

// 缓存域名信息
func (c *cloudflare) cache() error {
	resp, err := c.listZone()
	if err != nil {
		return err
	}
	for _, v := range resp.Result {
		c.domains[v.Name] = v.Id
	}
	return nil
}

func NewCloudflareDomain(token, domain string) (*cloudflareDomain, error) {
	cf, err := NewCloudFlare(token)
	if err != nil {
		return nil, err
	}
	return &cloudflareDomain{
		client: cf,
		domain: domain,
	}, nil
}

type cloudflareDomain struct {
	client *cloudflare
	domain string
}

func cf2rcds(prefix string, records ...*cloudflareRecord) []*ddns.Record {
	var res = make([]*ddns.Record, 0, len(records))
	for _, record := range records {
		res = append(res, &ddns.Record{
			Type:       record.Type,
			DomainName: record.Name[len(prefix):],
			RecordID:   record.Id,
			Value:      record.Content,
			Prefix:     prefix,
		})
	}
	return res
}

func rcd2cfs(records ...*ddns.Record) []*cloudflareRecord {
	var res = make([]*cloudflareRecord, 0, len(records))
	for _, record := range records {
		rcd := &cloudflareRecord{
			Id:      record.RecordID,
			Type:    record.Type,
			Content: record.Value,
		}
		if record.Prefix == "@" {
			rcd.Name = record.DomainName
		} else {
			rcd.Name = record.Prefix + "." + record.DomainName
		}
		res = append(res, rcd)
	}
	return res
}

func (c cloudflareDomain) GetByPrefix(prefix string) ([]*ddns.Record, error) {
	res, err := c.client.listRecords(c.domain, prefix)
	if err != nil {
		return nil, err
	}
	return cf2rcds(prefix, res.Result...), nil
}

func (c cloudflareDomain) CreateByRecords(rcds ...*ddns.Record) error {
	cfrcds := rcd2cfs(rcds...)
	for i, v := range cfrcds {
		err := c.client.createRecords(v)
		if err != nil {
			return err
		}
		*rcds[i] = *cf2rcds(rcds[i].Prefix, v)[0]
	}
	return nil
}

func (c cloudflareDomain) UpdateByRecords(rcds ...*ddns.Record) error {
	cfrcds := rcd2cfs(rcds...)
	for i, v := range cfrcds {
		err := c.client.updateRecords(v)
		if err != nil {
			return err
		}
		*rcds[i] = *cf2rcds(rcds[i].Prefix, v)[0]
	}
	return nil
}

func (c cloudflareDomain) DeleteByRecords(rcds ...*ddns.Record) error {
	cfrcds := rcd2cfs(rcds...)
	for _, v := range cfrcds {
		err := c.client.deleteRecords(v)
		if err != nil {
			return err
		}
	}
	return nil
}
