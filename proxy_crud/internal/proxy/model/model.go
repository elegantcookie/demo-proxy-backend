package model

import (
	"github.com/google/uuid"
	"time"
)

type Proxy struct {
	ID               string    `json:"id" structs:"id"`
	Ip               string    `json:"ip" structs:"ip"`
	Port             int       `json:"port" structs:"port"`
	ExternalIP       string    `json:"external_ip" structs:"external_ip"`
	Country          string    `json:"country" structs:"country"`
	OpenPorts        []int     `json:"open_ports" structs:"open_ports"`
	Active           bool      `json:"active" structs:"active"`
	Ping             int       `json:"ping" structs:"ping"`
	CreatedAt        time.Time `json:"created_at" structs:"created_at"`
	CheckedAt        time.Time `json:"checked_at" structs:"checked_at"`
	ValidAt          time.Time `json:"valid_at" structs:"valid_at"`
	BLCheck          int       `json:"bl_check" structs:"bl_check"`
	ProcessingStatus int       `json:"processing_status" structs:"processing_status"`
	ProxyGroupID     string    `json:"proxy_group_id" structs:"proxy_group_id"`
}

// NewProxy transfers DTO to proxy, generating uuid and time of creation
func NewProxy(dto CreateProxyDTO) Proxy {
	return Proxy{
		ID:           uuid.NewString(),
		Ip:           dto.Ip,
		Port:         dto.Port,
		ExternalIP:   dto.ExternalIP,
		Country:      dto.Country,
		CreatedAt:    time.Now(),
		ProxyGroupID: dto.ProxyGroupID,
	}
}

func NewProxyFromUpdateDTO(id string, dto UpdateProxyDTO) Proxy {
	return Proxy{
		ID:               id,
		Ip:               dto.Ip,
		Port:             dto.Port,
		ExternalIP:       dto.ExternalIP,
		Country:          dto.Country,
		OpenPorts:        dto.OpenPorts,
		Active:           dto.Active,
		Ping:             dto.Ping,
		CreatedAt:        dto.CreatedAt,
		CheckedAt:        dto.CheckedAt,
		ValidAt:          dto.ValidAt,
		BLCheck:          dto.BLCheck,
		ProcessingStatus: dto.ProcessingStatus,
		ProxyGroupID:     dto.ProxyGroupID,
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
