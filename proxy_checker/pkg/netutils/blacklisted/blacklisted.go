package blacklisted

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"net"
	"sync/atomic"
)

func Blacklisted(ip string) (bool, error) {
	sem := semaphore.NewWeighted(int64(20))
	ipv4, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		return false, err
	}
	hosts := DNSBlackLists
	items := make([]Item, 0)
	for _, host := range hosts {
		items = append(items, Item{
			IP:        ipv4.IP,
			Blacklist: fmt.Sprintf("%s.%s.", ReverseIP(ipv4.String()), host),
			Host:      host,
		})
	}
	var blacklisted uint64

	for _, i := range items {
		if err = sem.Acquire(context.Background(), 1); err != nil {
			return false, err
		}
		go func(item Item) {
			if b, _ := processIpCheck(sem, item); b {
				atomic.AddUint64(&blacklisted, 1)
			}
		}(i)
	}

	if err = sem.Acquire(context.Background(), int64(20)); err != nil {
		return false, err
	}

	if blacklisted > 0 {
		return true, nil
	}

	return false, nil
}
