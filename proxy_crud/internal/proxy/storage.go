package proxy

import "context"

type Storage interface {
	Insert(ctx context.Context, proxies []Proxy) error
	InsertOne(ctx context.Context, proxy Proxy) error
	FindById(ctx context.Context, id string) (Proxy, error)
	FindAll(ctx context.Context) ([]Proxy, error)
	Update(ctx context.Context, proxy Proxy) error
	Delete(ctx context.Context, id string) error
}
