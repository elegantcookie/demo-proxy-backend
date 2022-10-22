package model

import (
	"fmt"
	"strconv"
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
