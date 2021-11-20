package ip

import (
	"fmt"
	"testing"
)

func TestXdIpGetter(t *testing.T) {
	getter := XdIpGetter{}
	fmt.Println(getter.IP())
}
