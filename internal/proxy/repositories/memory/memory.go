package memory

import (
	"context"
	"fmt"
	"github.com/zhanbolat18/proxy/internal/proxy"
	"sync"
)

type memoryRepository struct {
	mainMap *sync.Map
}

func NewMemoryRepository() proxy.ProxyRepository {
	m := &sync.Map{}
	return &memoryRepository{mainMap: m}
}

func (m *memoryRepository) Store(_ context.Context, id string, object *proxy.Data) error {
	m.mainMap.Store(id, object)
	return nil
}

func (m *memoryRepository) GetAll(_ context.Context) map[string]*proxy.Data {
	res := make(map[string]*proxy.Data)
	m.mainMap.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(*proxy.Data)
		res[k] = v
		return true
	})
	return res
}

func (m *memoryRepository) GetById(_ context.Context, id string) (*proxy.Data, error) {
	data, ok := m.mainMap.Load(id)
	if !ok {
		return nil, fmt.Errorf("data with id %s not found", id)
	}
	return data.(*proxy.Data), nil
}

func (m *memoryRepository) initMap() {
	if m.mainMap == nil {
		m.mainMap = &sync.Map{}
	}
}
