package detector

import (
	"fmt"
	"testing"
)

func TestPublicIpGetter_IP(t *testing.T) {
	getter := Public{}
	fmt.Println(getter.IP())
}
