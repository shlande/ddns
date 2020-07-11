package ddns

type DNS interface {
	Update() error
}

type DomainInfo struct {
	Type       string
	DomainName string
	Prefix     string
}
