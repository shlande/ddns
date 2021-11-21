package detector

import (
	"fmt"
	"testing"
)

func TestGetIpFromDevice(t *testing.T) {
	getter, err := NewDeviceGetter("en0")
	if err != nil {
		panic(err)
	}
	ip, err := getter.IP()
	if err != nil {
		panic(err)
	}
	fmt.Println(ip)
	ip, err = getter.IPv6()
	if err != nil {
		panic(err)
	}
	fmt.Println(ip)
}
