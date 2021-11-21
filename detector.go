package ddns

type Detector interface {
	IP() ([]string, error)
}

type DetectorV6 interface {
	IPv6() ([]string, error)
}
