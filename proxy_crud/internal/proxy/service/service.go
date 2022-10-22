package service

import (
	"context"
	"proxy_crud/internal/proxy/db"
	"proxy_crud/internal/proxy/model"
	"proxy_crud/internal/proxy/pstorage"
	"proxy_crud/internal/proxy_group/pgstorage"
	"proxy_crud/pkg/api/filter"
	"proxy_crud/pkg/logging"
)

var _ Service = &service{}

type service struct {
	//cache   Cache
	ProxyGroupStorage pgstorage.Storage
	ProxyStorage      pstorage.Storage
	Logger            logging.Logger
}

func NewService(proxyStorage pstorage.Storage, proxyGroupStorage pgstorage.Storage, logger *logging.Logger) (Service, error) {
	return &service{
		ProxyStorage:      proxyStorage,
		ProxyGroupStorage: proxyGroupStorage,
		//cache:   Cache,
		Logger: *logger,
	}, nil
}

func (s service) createProxyDTOisValid(ctx context.Context, dto model.CreateProxyDTO) bool {
	_, err := s.ProxyGroupStorage.FindById(ctx, dto.ProxyGroupID)
	if err != nil {
		return false
	}
	return true
}

func (s service) validateProxyDTOs(ctx context.Context, dto []model.CreateProxyDTO) []model.CreateProxyDTO {
	validDTOs := make([]model.CreateProxyDTO, 0)
	for i := 0; i < len(dto); i++ {
		if s.createProxyDTOisValid(ctx, dto[i]) {
			validDTOs = append(validDTOs, dto[i])
		}
	}
	return validDTOs
}

func (s service) AddProxies(ctx context.Context, dto []model.CreateProxyDTO) error {
	validDTOs := s.validateProxyDTOs(ctx, dto)
	proxies := model.NewProxies(validDTOs)
	err := s.ProxyStorage.Insert(ctx, proxies)
	if err != nil {
		return err
	}
	return nil
}

func (s service) AddProxy(ctx context.Context, dto model.CreateProxyDTO) error {
	p := model.NewProxy(dto)
	err := s.ProxyStorage.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetAll(ctx context.Context, dto filter.Options) ([]model.Proxy, error) {
	options := db.NewPSQLFilterOptions(dto)
	proxies, err := s.ProxyStorage.FindAll(ctx, options)
	if err != nil {
		return nil, err
	}
	return proxies, nil
}

func (s service) GetById(ctx context.Context, id string) (model.Proxy, error) {
	proxy, err := s.ProxyStorage.FindById(ctx, id)
	if err != nil {
		return model.Proxy{}, err
	}
	return proxy, nil
}

func (s service) Update(ctx context.Context, id string, dto model.UpdateProxyDTO) error {
	pr := model.NewProxyFromUpdateDTO(id, dto)
	err := s.ProxyStorage.Update(ctx, pr)
	if err != nil {
		return err
	}
	return nil
}

func (s service) UpdateProxyStatus(ctx context.Context, id string, status int) error {
	proxy, err := s.ProxyStorage.FindById(ctx, id)
	if err != nil {
		return err
	}
	proxy.ProcessingStatus = status
	err = s.ProxyStorage.Update(ctx, proxy)
	if err != nil {
		return err
	}
	return nil
}

func (s service) DeleteAll(ctx context.Context) error {
	err := s.ProxyStorage.DeleteAll(ctx)
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
	Update(ctx context.Context, id string, dto model.UpdateProxyDTO) error
	//PartialUpdate(ctx context.Context, )
	UpdateProxyStatus(ctx context.Context, id string, status int) error
	//Delete(ctx context.Context, uuid string) error
	DeleteAll(ctx context.Context) error
}
