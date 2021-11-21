package ddns

type DomainBuilder func(domain string, id, key string) (Domain, error)

type Domain interface {
	GetByPrefix(prefix string) ([]*Record, error)
	CreateByRecords(rcds ...*Record) error
	UpdateByRecords(rcd ...*Record) error
	DeleteByRecords(rcd ...*Record) error
}

func WithCache(domain Domain) *cache {
	return &cache{
		Domain: domain,
		cache:  map[string]*Record{},
	}
}

// cache可以用来缓存一部分的信息
// 这部分信息在删除和更新的时候可以减少查找的次数
type cache struct {
	Domain
	// key：recordId
	cache map[string]*Record
}

func (d *cache) GetByPrefix(prefix string) ([]*Record, error) {
	rcds, err := d.Domain.GetByPrefix(prefix)
	if err != nil {
		return nil, err
	}
	d.save(rcds...)
	return rcds, nil
}

func (d cache) CreateByRecords(rcds ...*Record) error {
	if len(rcds) == 0 {
		return nil
	}
	err := d.Domain.CreateByRecords(rcds...)
	if err != nil {
		return err
	}
	d.save(rcds...)
	return nil
}

func (d cache) save(rcds ...*Record) {
	for _, v := range rcds {
		d.cache[v.RecordID] = v
	}
}

// UpdateByRecords 注意，update不会更新参数中的信息
func (d cache) UpdateByRecords(rcd ...*Record) error {
	if len(rcd) == 0 {
		return nil
	}
	for _, v := range rcd {
		d.cache[v.RecordID].Merge(v)
	}
	return d.Domain.UpdateByRecords(rcd...)
}

func (d cache) DeleteByRecords(rcd ...*Record) error {
	if len(rcd) == 0 {
		return nil
	}
	err := d.Domain.DeleteByRecords(rcd...)
	if err != nil {
		return err
	}
	for _, v := range rcd {
		delete(d.cache, v.RecordID)
	}
	return nil
}
