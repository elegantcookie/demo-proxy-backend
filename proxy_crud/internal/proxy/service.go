package proxy

import (
	"context"
	"proxy_crud/pkg/logging"
)

var _ Service = &service{}

type service struct {
	//cache   Cache
	storage Storage
	logger  logging.Logger
}

func (s service) AddProxies(ctx context.Context, dto []CreateProxyDTO) error {
	proxies := NewProxies(dto)
	err := s.storage.Insert(ctx, proxies)
	if err != nil {
		return err
	}
	return nil
}

func (s service) AddProxy(ctx context.Context, dto CreateProxyDTO) error {
	p := NewProxy(dto)
	err := s.storage.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetAll(ctx context.Context) ([]Proxy, error) {
	proxies, err := s.storage.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return proxies, nil
}

func (s service) GetById(ctx context.Context, id string) (Proxy, error) {
	proxy, err := s.storage.FindById(ctx, id)
	if err != nil {
		return Proxy{}, err
	}
	return proxy, nil
}

type Service interface {
	AddProxies(ctx context.Context, dto []CreateProxyDTO) error
	AddProxy(ctx context.Context, dto CreateProxyDTO) error
	GetAll(ctx context.Context) ([]Proxy, error)
	GetById(ctx context.Context, id string) (Proxy, error)
	//Update(ctx context.Context, dto UpdateProxyDTO) error
	//Delete(ctx context.Context, uuid string) error
}
