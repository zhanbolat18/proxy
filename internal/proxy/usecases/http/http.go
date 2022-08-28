package http

import (
	"bytes"
	"context"
	"errors"
	"github.com/zhanbolat18/proxy/internal/proxy"
	"github.com/zhanbolat18/proxy/internal/proxy/entities"
	"github.com/zhanbolat18/proxy/pkg/doer"
	"github.com/zhanbolat18/proxy/pkg/generator"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type httpUsecase struct {
	doer        doer.Doer
	idGenerator generator.Generator
	repo        proxy.ProxyRepository
}

func NewHttpUsecase(doer doer.Doer, idGenerator generator.Generator, repo proxy.ProxyRepository) proxy.ProxyUsecase {
	return &httpUsecase{doer: doer, idGenerator: idGenerator, repo: repo}
}

func (h *httpUsecase) Proxy(ctx context.Context, request entities.Request) (*entities.Response, error) {
	err := h.validateUrl(request.Url)
	if err != nil {
		return nil, err
	}
	req, err := h.makeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	res, err := h.doer.Do(req)
	if err != nil {
		return nil, err
	}
	response, err := h.handleResponse(res)
	if err != nil {
		return nil, err
	}
	err = h.repo.Store(ctx, response.Id, &proxy.Data{
		Req: &request,
		Res: response,
	})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (h *httpUsecase) handleResponse(res *http.Response) (*entities.Response, error) {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	id, err := h.idGenerator.Id()
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	for s, v := range res.Header {
		headers[s] = strings.Join(v, ";")
	}
	length := res.ContentLength
	if length < 0 && len(body) > 0 {
		length = int64(len(body))
	}

	return &entities.Response{
		Id:      id,
		Status:  res.StatusCode,
		Headers: headers,
		Length:  length,
		Body:    body,
	}, nil
}

func (h *httpUsecase) makeRequest(ctx context.Context, request entities.Request) (*http.Request, error) {
	var body io.Reader = nil
	if len(request.Body) > 0 {
		body = bytes.NewReader(request.Body)
	}
	req, err := http.NewRequestWithContext(ctx, request.Method, request.Url, body)
	if err != nil {
		return nil, err
	}
	if len(request.Headers) > 0 {
		for header, value := range request.Headers {
			req.Header.Set(header, value)
		}
	}
	return req, nil
}

func (h *httpUsecase) validateUrl(checkUrl string) error {
	u, err := url.Parse(checkUrl)
	if err != nil {
		return err
	}
	if u.Scheme == "" || u.Host == "" {
		return errors.New("invalid url")
	}
	return nil
}
