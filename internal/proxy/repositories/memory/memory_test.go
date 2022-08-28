package memory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zhanbolat18/proxy/internal/proxy"
	"sync"
	"testing"
)

var repo = &memoryRepository{}
var ctx = context.Background()

func TestMemoryRepository_Store(t *testing.T) {
	m := &sync.Map{}
	repo.mainMap = m
	d := &proxy.Data{Req: nil, Res: nil}
	err := repo.Store(ctx, "id", d)
	assert.Nil(t, err)
	data, ok := m.Load("id")
	assert.True(t, ok)
	assert.NotNil(t, data)
	assert.Equal(t, d, data)
}
