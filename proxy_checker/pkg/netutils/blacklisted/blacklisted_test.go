package blacklisted

import (
	"fmt"
	"testing"
)

func TestBlacklisted(t *testing.T) {
	blacklisted, err := Blacklisted("1.2.104.208")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(blacklisted)
}
