package ddns

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

func RunIpv6Worker(ctx context.Context, getter DetectorV6, binder *Binder, ttl int64) {
	runTTl(ctx, ttl, func() error {
		addr, err := getter.IPv6()
		if err != nil {
			return err
		}
		return binder.BindV6(addr)
	})
}

func RunIpWorker(ctx context.Context, ip Detector, binder *Binder, ttl int64) {
	runTTl(ctx, ttl, func() error {
		addr, err := ip.IP()
		if err != nil {
			return err
		}
		return binder.Bind(addr)
	})
}

func runTTl(ctx context.Context, ttl int64, worker func() error) {
	ticker := time.NewTimer(time.Duration(ttl * int64(time.Second)))
	defer ticker.Stop()
	for {
		err := worker()
		if err != nil {
			logrus.Error(err)
		}
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				break
			}
		}
	}
}
