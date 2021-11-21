# DDNS

go语言实现的ddns，能将ip自动绑定到指定域名。

## 使用方法

### 安装

```
go 
```

### DNS

支持的dns服务提供方有

1. 阿里云
2. Dnspod（腾讯云）
3. cloudflare（todo）

### 获取ip

通过detect参数来指定。

1. 从网卡获取地址（支持ipv4/ipv6）：
   
   ip 参数传入 `device=<match>`，其中match支持正则表达式，如果匹配到多个网卡，那么会把所有的网络地址都绑定到域名上

   例如： 
   
   ```
   --detect=device=ens.* # 获取所有以ens开头的网卡的ip
   --detect=device=eth0 # 获取网卡名为eth0的ip
   ```
   
2. 获取公网ip（仅支持ipv4）：

   使用ip-addr的公开接口查找发起请求方的ip地址。

3. xd内网（仅支持ipv4）

   获取校园网内网ip，并绑定到域名上。

### 例子

每过十秒检查一次校园网ip，并绑定到test.shlande.top中，dns服务提供方是dnspod

```
ddns --provider=dnspod --domain=colaha.tech --prefix=test --detect=xd --type=ip --secret-id=<id_here> --secret-key=<key-here>
```

## TODO

1. 支持配置文件绑定
2. 允许多个域名同时绑定
3. detect模仿calico，支持连通性测试后自动选择
