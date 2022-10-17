package parser

import (
	"fmt"
	"strings"
)

type ParsedFields struct {
	Ip         string
	Port       string
	ExternalIp string
	Country    string
}

// SplitText returns non-empty strings split by '\n' symbol in text
func SplitText(text string) []string {
	return strings.Split(text, "\n")
}

// ParseLine parses and validates proxy line
func ParseLine(proxyLine string) (pf ParsedFields, err error) {
	strings.TrimSuffix(proxyLine, "\r")
	if len(proxyLine) == 0 {
		return pf, fmt.Errorf("proxy string is empty")
	}
	fields := strings.Split(proxyLine, ";")
	if len(fields) < 4 {
		return pf, fmt.Errorf("proxy string contains less than 4 fields")
	}
	ipPort := strings.Split(fields[0], ":")
	if len(ipPort) != 2 {
		return pf, fmt.Errorf("wrong ip:port format")
	}
	pf.Ip = ipPort[0]
	pf.Port = ipPort[1]
	pf.ExternalIp = fields[1]
	pf.Country = fields[3]
	return pf, nil
}
