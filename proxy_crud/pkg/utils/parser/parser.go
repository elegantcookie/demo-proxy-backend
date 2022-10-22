package parser

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type ParsedFields struct {
	Ip           string
	Port         string
	ExternalIp   string
	Country      string
	ProxyGroupID string
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

	fields := strings.Split(proxyLine, "|")
	if len(fields) != 2 {
		return pf, fmt.Errorf("wrong proxy string format")
	}

	proxyFields := strings.Split(fields[0], ";")
	if len(proxyFields) < 4 {
		return pf, fmt.Errorf("proxy fields string contains less than 4 proxyFields")
	}
	ipPort := strings.Split(proxyFields[0], ":")
	if len(ipPort) != 2 {
		return pf, fmt.Errorf("wrong ip:port format")
	}

	pf.Ip = ipPort[0]
	pf.Port = ipPort[1]
	pf.ExternalIp = proxyFields[1]
	pf.Country = proxyFields[3]

	proxyGroupId := fields[1]
	if len(proxyGroupId) == 0 {
		return pf, nil
	}
	proxyGroupId = strings.TrimSuffix(proxyGroupId, "\r")
	if proxyGroupId == "" {
		return pf, fmt.Errorf("proxy group id is empty")
	}
	_, err = uuid.Parse(proxyGroupId)
	if err != nil {
		return pf, fmt.Errorf("proxy group id has invalid format")
	}
	pf.ProxyGroupID = proxyGroupId

	return pf, nil
}
