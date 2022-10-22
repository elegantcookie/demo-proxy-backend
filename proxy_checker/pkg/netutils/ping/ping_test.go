package ping

import (
	"fmt"
	"testing"
)

func TestPing(t *testing.T) {
	t.Run("existing proxy ping test", func(t *testing.T) {
		ping, err := Ping("141.101.123.228")
		if err != nil {
			fmt.Printf("failed to ping: %v\n", err)
		}
		fmt.Println(ping)
	})

	t.Run("existing proxy ping test", func(t *testing.T) {
		ping, err := Ping("45.8.105.232")
		if err != nil {
			fmt.Printf("failed to ping: %v\n", err)
		}
		fmt.Println(ping)
	})
}
