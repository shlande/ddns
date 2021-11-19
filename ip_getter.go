package ddns

type IpGetter interface {
	IP() ([]string, error)
}

type IPv6Getter interface {
	IPv6() ([]string, error)
}
