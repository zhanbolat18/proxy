package http_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zhanbolat18/proxy/internal/proxy/entities"
	"github.com/zhanbolat18/proxy/internal/proxy/repositories/memory"
	http2 "github.com/zhanbolat18/proxy/internal/proxy/usecases/http"
	"github.com/zhanbolat18/proxy/pkg/generator/uuid"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// box testing
type mockDoer struct {
	mock.Mock
}

func (m *mockDoer) Do(r *http.Request) (*http.Response, error) {
	args := m.Called(r.URL.String(), r.Method)
	var res *http.Response = nil
	if args.Get(0) != nil {
		res = args.Get(0).(*http.Response)
	}
	return res, args.Error(1)
}

var doer = &mockDoer{}
var idGen = uuid.NewUUIDGenerator()
var repo = memory.NewMemoryRepository()
var uc = http2.NewHttpUsecase(doer, idGen, repo)
var ctx = context.Background()

func TestHttpUsecase_ProxySuccess(t *testing.T) {
	req := entities.Request{
		Method:  http.MethodGet,
		Url:     "http://example.com",
		Headers: nil,
		Body:    nil,
	}
	bodyText := "some body text"
	doer.On("Do", req.Url, req.Method).Return(&http.Response{
		Status:        http.StatusText(http.StatusOK),
		StatusCode:    http.StatusOK,
		Body:          ioutil.NopCloser(strings.NewReader(bodyText)),
		ContentLength: int64(len(bodyText)),
	}, nil)
	res, err := uc.Proxy(ctx, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Id)
	assert.Equal(t, res.Status, http.StatusOK)
	assert.Equal(t, string(res.Body), bodyText)
	assert.Equal(t, res.Length, int64(len(bodyText)))

	data, err := repo.GetById(ctx, res.Id)
	assert.Nil(t, err)
	assert.Equal(t, *data.Req, req)
	assert.Equal(t, data.Res, res)
}

func TestHttpUsecase_ProxyInvalidUrl(t *testing.T) {
	req := entities.Request{
		Method:  http.MethodGet,
		Url:     "adsfdsf",
		Headers: nil,
		Body:    nil,
	}

	res, err := uc.Proxy(ctx, req)
	doer.AssertNotCalled(t, "Do")
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestHttpUsecase_ProxyInvalidMethod(t *testing.T) {
	req := entities.Request{
		Method:  "invalid method",
		Url:     "http://example.com",
		Headers: nil,
		Body:    nil,
	}

	res, err := uc.Proxy(ctx, req)
	doer.AssertNotCalled(t, "Do")
	assert.NotNil(t, err)
	assert.Nil(t, res)
}
