package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"proxy_checker/pkg/logging"
	"proxy_checker/pkg/netutils/blacklisted"
	"proxy_checker/pkg/netutils/ping"
	"strings"
	"time"
)

var _ Service = &service{}

const (
	statusToBeChecked         = 0
	statusCheckWasInterrupted = 2
	proxyCRUDUpdateURL        = "http://proxy_crud:10000/api/proxy_crud/v1/proxy/upd/id/"
)

type service struct {
	Logger logging.Logger
}

func (s service) FetchChanges(ctx context.Context, proxy Proxy) error {
	bytes, err := json.Marshal(proxy)
	if err != nil {
		return err
	}
	body := io.NopCloser(strings.NewReader(string(bytes)))
	url := proxyCRUDUpdateURL + proxy.ID
	request, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return err
	}

	// TODO: configure client
	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		bytes, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("wrong status code: %d, body: %s", response.StatusCode, string(bytes))
	}
	return nil
}

func (s service) Check(ctx context.Context, pr Proxy) (Proxy, error) {
	pr.ProcessingStatus = statusToBeChecked

	p, err := ping.Ping(pr.Ip)
	if err != nil {
		pr.Active = false
		return pr, err
	}
	var blStatus int
	if bl, err := blacklisted.Blacklisted(pr.Ip); err == nil {
		if bl {
			blStatus = 2
		} else {
			blStatus = 1
		}
	}
	pr.Ping = int(p.Milliseconds())
	if !pr.CheckedAt.IsZero() {
		now := time.Now()
		pr.BLCheck = blStatus
		pr.CheckedAt = now
		pr.ValidAt = now
		pr.Active = true
		return pr, nil
	}

	now := time.Now()
	pr.BLCheck = blStatus
	pr.CheckedAt = now
	pr.ValidAt = now
	pr.Active = true

	// TODO: check country
	// TODO: check external ip

	return pr, nil
}

func NewService(logger *logging.Logger) (Service, error) {
	return &service{
		Logger: *logger,
	}, nil
}

type Service interface {
	FetchChanges(ctx context.Context, proxy Proxy) error
	Check(ctx context.Context, pr Proxy) (Proxy, error)
}
