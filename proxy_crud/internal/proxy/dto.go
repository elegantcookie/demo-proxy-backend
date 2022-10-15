package proxy

import "time"

type CreateProxyDTO struct {
	Ip         string `json:"ip"`
	Port       int    `json:"port"`
	ExternalIP string `json:"external_ip"`
	Country    string `json:"country"`
}

func NewProxy(dto CreateProxyDTO) Proxy {
	return Proxy{
		Ip:         dto.Ip,
		Port:       dto.Port,
		ExternalIP: dto.ExternalIP,
		Country:    dto.Country,
		CreatedAt:  time.Time{},
	}
}

func NewProxies(dto []CreateProxyDTO) []Proxy {
	proxies := make([]Proxy, len(dto))
	for i := 0; i < len(dto); i++ {
		p := NewProxy(dto[i])
		proxies[i] = p
	}
	return proxies
}
