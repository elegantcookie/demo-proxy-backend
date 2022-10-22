package service

import (
	"context"
	"proxy_crud/internal/proxy_group/db"
	"proxy_crud/internal/proxy_group/model"
	"proxy_crud/internal/proxy_group/pgstorage"
	"proxy_crud/pkg/api/filter"
	"proxy_crud/pkg/logging"
)

var _ Service = &service{}

type service struct {
	//cache   Cache
	storage pgstorage.Storage
	logger  logging.Logger
}

func NewService(proxyGroupStorage pgstorage.Storage, logger *logging.Logger) (Service, error) {
	return &service{
		storage: proxyGroupStorage,
		logger:  *logger,
	}, nil
}

func (s service) AddProxyGroup(ctx context.Context, dto model.CreateProxyGroupDTO) error {
	p := model.NewProxyGroup(dto)
	err := s.storage.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetAll(ctx context.Context, dto filter.Options) ([]model.ProxyGroup, error) {
	options := db.NewPSQLFilterOptions(dto)
	proxies, err := s.storage.FindAll(ctx, options)
	if err != nil {
		return nil, err
	}
	return proxies, nil
}

func (s service) GetById(ctx context.Context, id string) (model.ProxyGroup, error) {
	proxy, err := s.storage.FindById(ctx, id)
	if err != nil {
		return model.ProxyGroup{}, err
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
	AddProxyGroup(ctx context.Context, dto model.CreateProxyGroupDTO) error
	GetAll(ctx context.Context, options filter.Options) ([]model.ProxyGroup, error)
	GetById(ctx context.Context, id string) (model.ProxyGroup, error)
	//Update(ctx context.Context, dto UpdateProxyGroupDTO) error
	//Delete(ctx context.Context, uuid string) error
	DeleteAll(ctx context.Context) error
}
