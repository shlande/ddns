package ddns

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

func runIpv6Worker(ctx context.Context, getter IPv6Getter, v6 DnsV6, ttl int64) {
	runTTl(ctx, ttl, func() error {
		addr, err := getter.IPv6()
		if err != nil {
			return err
		}
		return v6.BindV6(addr)
	})
}

func runIpWorker(ctx context.Context, ip IpGetter, dns Dns, ttl int64) {
	runTTl(ctx, ttl, func() error {
		addr, err := ip.IP()
		if err != nil {
			return err
		}
		return dns.Bind(addr)
	})
}

func runTTl(ctx context.Context, ttl int64, worker func() error) {
	ticker := time.NewTimer(time.Duration(ttl * int64(time.Second)))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := worker()
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}
