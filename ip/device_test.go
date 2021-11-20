package ip

import (
	"fmt"
	"testing"
)

func TestGetIpFromDevice(t *testing.T) {
	getter := NewDeviceGetter("en0")
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
