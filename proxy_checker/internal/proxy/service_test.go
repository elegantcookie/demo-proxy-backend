package proxy

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	/*logger := logging.GetLogger("info")
	service, _ := NewService(&logger)*/
	t.Run("proxy ping", func(t *testing.T) {
		link := "http://212.16.19.1"
		res, err := http.Get(link)
		if err != nil {
			fmt.Println(err)
			return
		}

		body1, _ := io.ReadAll(res.Body)

		fmt.Println(string(body1))
	})
}
