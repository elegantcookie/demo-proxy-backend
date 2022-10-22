package pgstorage

import (
	"context"
	"proxy_crud/internal/proxy_group/model"
)

type Storage interface {
	InsertOne(ctx context.Context, proxy model.ProxyGroup) error
	FindById(ctx context.Context, id string) (model.ProxyGroup, error)
	FindAll(ctx context.Context, options IOptions) ([]model.ProxyGroup, error)
	Update(ctx context.Context, proxy model.ProxyGroup) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}

type IOptions interface {
	GetOrderBy() string
	GetLimit() uint64
	GetOffset() uint64
	MapOptions() Sqlizer
}

type Sqlizer interface {
	ToSql() (string, []interface{}, error)
}
