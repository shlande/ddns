package dns

type Domain interface {
	GetByPrefix(prefix string) ([]*Record, error)
	CreateByRecords(rcds ...*Record) error
	UpdateByRecords(rcd ...*Record) error
	DeleteByRecords(rcd ...*Record) error
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
	for _, v := range rcd {
		mergeRecord(d.cache[v.RecordID], v)
	}
	return d.Domain.UpdateByRecords(rcd...)
}

func (d cache) DeleteByRecords(rcd ...*Record) error {
	err := d.Domain.DeleteByRecords(rcd...)
	if err != nil {
		return err
	}
	for _, v := range rcd {
		delete(d.cache, v.RecordID)
	}
	return nil
}
