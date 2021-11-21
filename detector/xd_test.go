package detector

import (
	"fmt"
	"testing"
)

func TestXdIpGetter(t *testing.T) {
	getter := Xd{}
	fmt.Println(getter.IP())
}
