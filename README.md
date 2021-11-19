## DDNS
go语言实现的ddns，能将公网ip自动绑定到指定域名，支持：

1. ipv6
2. ipv4
3. xd校园网（ipv4）

```go
package main
import "github.com/shland/ddns"

func main() {
    d := ddns.NewAliDNS("key","secret",ddns.DomainInfo{"name","prefix"})
    d.Update()
}
```
TODO:
1. 支持从网卡中获取ip

   这里的规则仿照calico的规则：

    1. 支持通过名称来匹配
    2. 支持通过是否能够ping同某个地址来匹配

2. 重构一下：

   分为：IPGetter，DnsProvider。

   IPGetter负责提供ip信息，DnsProvider负责把指定信息绑定好。

