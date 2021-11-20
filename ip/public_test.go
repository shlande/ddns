package ip

import (
	"fmt"
	"testing"
)

func TestPublicIpGetter_IP(t *testing.T) {
	getter := PublicIpGetter{}
	fmt.Println(getter.IP())
}
