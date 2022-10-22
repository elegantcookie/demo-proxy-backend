package proxy

import "time"

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
	ProxyGroupID     string    `json:"proxy_group_id"`
}

/*func NewProxy(dto CheckProxyDTO) Proxy {
	return Proxy{
		ID:           dto.ID,
		Ip:           dto.Ip,
		Port:         dto.Port,
		ExternalIP:   dto.ExternalIP,
		Country:      dto.Country,
		CreatedAt:    dto.CreatedAt,
		CheckedAt:    time.Now(),
		ValidAt:      dto.ValidAt,
		ProxyGroupID: dto.ProxyGroupID,
	}
}
*/
