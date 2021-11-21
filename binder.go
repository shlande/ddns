package ddns

func NewBinder(prefix string, domain Domain) *Binder {
	return &Binder{
		domain: domain,
		prefix: prefix,
	}
}

type Binder struct {
	domain Domain
	prefix string
}

func (a *Binder) bind(addrs []string, tp string) error {
	rcds, err := a.domain.GetByPrefix(a.prefix)
	crt, del, upd, err := findBestSolution(rcds, tp, addrs)
	if err != nil {
		return err
	}
	err = a.domain.DeleteByRecords(del...)
	if err != nil {
		return err
	}
	err = a.domain.UpdateByRecords(upd...)
	if err != nil {
		return err
	}
	err = a.domain.CreateByRecords(crt...)
	if err != nil {
		return err
	}
	return nil
}

func (a *Binder) BindV6(addrs []string) error {
	return a.bind(addrs, "AAAA")
}

func (a *Binder) Bind(addrs []string) error {
	return a.bind(addrs, "A")
}
