package ddns

type DNS interface {
	Update() error
}

type DomainInfo struct {
	DomainName string
	Prefix     string
}


