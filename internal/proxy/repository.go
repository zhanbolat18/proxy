package proxy

import (
	"context"
	"github.com/zhanbolat18/proxy/internal/proxy/entities"
)

type Data struct {
	Req *entities.Request
	Res *entities.Response
}

type ProxyRepository interface {
	Store(ctx context.Context, id string, object *Data) error
	GetAll(ctx context.Context) map[string]*Data
	GetById(ctx context.Context, id string) (*Data, error)
}
