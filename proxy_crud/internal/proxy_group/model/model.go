package model

import (
	"github.com/google/uuid"
)

type ProxyGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewProxyGroup transfers DTO to proxy, generating uuid and time of creation
func NewProxyGroup(dto CreateProxyGroupDTO) ProxyGroup {
	return ProxyGroup{
		ID:   uuid.NewString(),
		Name: dto.Name,
	}
}
