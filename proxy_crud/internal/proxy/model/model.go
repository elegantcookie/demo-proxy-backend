package model

import (
	"github.com/google/uuid"
	"time"
)

type Proxy struct {
	ID               string    `json:"id"`
	Ip               string    `json:"ip"`
	Port             int       `json:"port"`
	ExternalIP       string    `json:"external_ip"`
	Country          string    `json:"country"`
	OpenPorts        []int     `json:"open_ports"`
	Active           bool      `json:"active"`
	Ping             int       `json:"ping"`
	CreatedAt        time.Time `json:"created_at"`
	CheckedAt        time.Time `json:"checked_at"`
	ValidAt          time.Time `json:"valid_at"`
	BLCheck          int       `json:"bl_check"`
	ProcessingStatus int       `json:"processing_status"`
}

// NewProxy transfers DTO to proxy, generating uuid and time of creation
func NewProxy(dto CreateProxyDTO) Proxy {
	return Proxy{
		ID:         uuid.NewString(),
		Ip:         dto.Ip,
		Port:       dto.Port,
		ExternalIP: dto.ExternalIP,
		Country:    dto.Country,
		CreatedAt:  time.Now(),
	}
}

// NewProxies transfers DTOs to proxies
func NewProxies(dto []CreateProxyDTO) []Proxy {
	proxies := make([]Proxy, len(dto))
	for i := 0; i < len(dto); i++ {
		proxies[i] = NewProxy(dto[i])
	}
	return proxies
}
