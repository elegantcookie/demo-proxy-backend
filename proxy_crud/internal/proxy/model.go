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
}
