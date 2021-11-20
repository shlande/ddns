package ddns

type Dns interface {
	Bind(addrs []string) error
}

type DnsV6 interface {
	BindV6(addrs []string) error
}
