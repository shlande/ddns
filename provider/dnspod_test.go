package provider

import (
	"fmt"
	"testing"
)

func TestDnsPodDomain_GetByPrefix(t *testing.T) {
	domain, err := NewDnsPodDomain("colaha.tech", SecretId, SecretKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(domain.GetByPrefix("test"))
}
