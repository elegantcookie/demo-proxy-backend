package pstorage

import (
	"context"
	"proxy_crud/internal/proxy/model"
)

type Storage interface {
	Insert(ctx context.Context, proxies []model.Proxy) error
	InsertOne(ctx context.Context, proxy model.Proxy) error
	FindById(ctx context.Context, id string) (model.Proxy, error)
	FindAll(ctx context.Context, options IFilterOptions) ([]model.Proxy, error)
	Update(ctx context.Context, proxy model.Proxy) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}

type IFilterOptions interface {
	GetOrderBy() string
}
