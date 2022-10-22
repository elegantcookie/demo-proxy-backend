package proxy

import "time"

type CheckProxyDTO struct {
	ID           string    `json:"id"`
	Ip           string    `json:"ip"`
	Port         int       `json:"port"`
	ExternalIP   string    `json:"external_ip"`
	Country      string    `json:"country"`
	Ping         int       `json:"ping"`
	CreatedAt    time.Time `json:"created_at"`
	ValidAt      time.Time `json:"valid_at"`
	ProxyGroupID string    `json:"proxy_group_id"`
}
