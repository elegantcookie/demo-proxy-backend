package model

import (
	"fmt"
	"strconv"
	"time"
)

type CreateProxyDTO struct {
	Ip           string `json:"ip"`
	Port         int    `json:"port"`
	ExternalIP   string `json:"external_ip"`
	Country      string `json:"country"`
	ProxyGroupID string `json:"proxy_group_id"`
}

func NewCreateProxyDTO(ip, port, externalIp, country, proxyGroupID string) (CreateProxyDTO, error) {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return CreateProxyDTO{}, fmt.Errorf("failed to convert port to int: %v", err)
	}
	return CreateProxyDTO{
		Ip:           ip,
		Port:         intPort,
		ExternalIP:   externalIp,
		Country:      country,
		ProxyGroupID: proxyGroupID,
	}, nil
}

type UpdateProxyDTO struct {
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
	ProxyGroupID     string    `json:"proxy_group_id"`
}

type UpdateProxyStatusDTO struct {
	ID               string `json:"id"`
	ProcessingStatus int    `json:"processing_status"`
}
