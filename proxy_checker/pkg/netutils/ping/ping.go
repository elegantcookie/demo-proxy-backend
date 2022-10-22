package ping

import (
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
)

func Ping(ip string) (rt time.Duration, err error) {
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		return rt, err
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		rt = rtt
	}
	err = p.Run()
	if err != nil {
		return rt, err
	}
	if err = p.Err(); err != nil {
		return rt, err
	}
	return rt, nil
}
