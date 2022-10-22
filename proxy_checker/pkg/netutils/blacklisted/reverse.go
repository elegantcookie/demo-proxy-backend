package blacklisted

import (
	"net"
	"strings"
)

func ReverseIP(IP string) string {
	var StringSplitIP []string

	if net.ParseIP(IP).To4() != nil { // Check for an IPv4 address
		StringSplitIP = strings.Split(IP, ".") // Split into 4 groups
		for x, y := 0, len(StringSplitIP)-1; x < y; x, y = x+1, y-1 {
			StringSplitIP[x], StringSplitIP[y] = StringSplitIP[y], StringSplitIP[x] // Reverse the groups
		}
	} else {
		StringSplitIP = strings.Split(IP, ":") // Split into however many groups

		/* Due to IPv6 lookups being different than IPv4 we have an extra check here
		We have to expand the :: and do 0-padding if there are less than 4 digits */
		for key := range StringSplitIP {
			if len(StringSplitIP[key]) == 0 { // Found the ::
				StringSplitIP[key] = strings.Repeat("0000", 8-strings.Count(IP, ":"))
			} else if len(StringSplitIP[key]) < 4 { // 0-padding needed
				StringSplitIP[key] = strings.Repeat("0", 4-len(StringSplitIP[key])) + StringSplitIP[key]
			}
		}

		// We have to join what we have and split it again to get all the letters split individually
		StringSplitIP = strings.Split(strings.Join(StringSplitIP, ""), "")

		for x, y := 0, len(StringSplitIP)-1; x < y; x, y = x+1, y-1 {
			StringSplitIP[x], StringSplitIP[y] = StringSplitIP[y], StringSplitIP[x]
		}
	}

	return strings.Join(StringSplitIP, ".") // Return the IP.
}
