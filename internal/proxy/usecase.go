package proxy

import (
	"context"
	"github.com/zhanbolat18/proxy/internal/proxy/entities"
)

type ProxyUsecase interface {
	Proxy(ctx context.Context, request entities.Request) (*entities.Response, error)
}
