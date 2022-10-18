package service

import (
	"context"
	"proxy_crud/internal/proxy/db"
	"proxy_crud/internal/proxy/model"
	"proxy_crud/internal/proxy/pstorage"
	"proxy_crud/pkg/api/filter"
	"proxy_crud/pkg/logging"
)

var _ Service = &service{}

type service struct {
	//cache   Cache
	storage pstorage.Storage
	logger  logging.Logger
}

func NewService(ProxyStorage pstorage.Storage, logger *logging.Logger) (Service, error) {
	return &service{
		storage: ProxyStorage,
		//cache:   Cache,
		logger: *logger,
	}, nil
}

func (s service) AddProxies(ctx context.Context, dto []model.CreateProxyDTO) error {
	proxies := model.NewProxies(dto)
	err := s.storage.Insert(ctx, proxies)
	if err != nil {
		return err
	}
	return nil
}

func (s service) AddProxy(ctx context.Context, dto model.CreateProxyDTO) error {
	p := model.NewProxy(dto)
	err := s.storage.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetAll(ctx context.Context, dto filter.Options) ([]model.Proxy, error) {
	options := db.NewPSQLFilterOptions(dto)
	proxies, err := s.storage.FindAll(ctx, options)
	if err != nil {
		return nil, err
	}
	return proxies, nil
}

func (s service) GetById(ctx context.Context, id string) (model.Proxy, error) {
	proxy, err := s.storage.FindById(ctx, id)
	if err != nil {
		return model.Proxy{}, err
	}
	return proxy, nil
}

func (s service) DeleteAll(ctx context.Context) error {
	err := s.storage.DeleteAll(ctx)
	if err != nil {
		return err
	}
	return nil
}

type Service interface {
	AddProxies(ctx context.Context, dto []model.CreateProxyDTO) error
	AddProxy(ctx context.Context, dto model.CreateProxyDTO) error
	GetAll(ctx context.Context, options filter.Options) ([]model.Proxy, error)
	GetById(ctx context.Context, id string) (model.Proxy, error)
	//Update(ctx context.Context, dto UpdateProxyDTO) error
	//Delete(ctx context.Context, uuid string) error
	DeleteAll(ctx context.Context) error
}
