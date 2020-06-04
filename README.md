## DDNS
go语言实现的ddns，能将公网ip自动绑定到指定域名
```go
package main
import "github.com/shland/ddns"

func main() {
    d := ddns.NewAliDNS("key","secret",ddns.DomainInfo{"name","prefix"})
    d.Update()
}
```