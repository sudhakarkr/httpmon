package httpmon

import (
	"bytes"
	"io"
	"net/http"
)

type HttpClient interface {
	Run(method HttpMethod, url string) (HttpResponse, error)
}

type HttpClientFactory func(out TimeOut) HttpClient

type HttpResponse interface {
	io.Closer
	StatusCode() HttpStatus
}

////

type defaultHttpClient struct {
	delegate http.Client
}

func (dhc *defaultHttpClient) Run(method HttpMethod, url string) (HttpResponse, error) {
	request, err := http.NewRequest(string(method), url, bytes.NewReader(make([]byte, 0)))
	if err != nil {
		return nil, err
	}
	response, err := dhc.delegate.Do(request)
	if err != nil {
		return nil, err
	}
	return &defaultHttpResponse{
		body:       response.Body,
		statusCode: response.StatusCode,
	}, nil
}

var defaultHttpClientFactory HttpClientFactory = func(out TimeOut) HttpClient {
	return &defaultHttpClient{delegate: http.Client{
		Timeout: out.ToDuration(),
	}}
}

type defaultHttpResponse struct {
	body       io.ReadCloser
	statusCode int
}

func (dhr *defaultHttpResponse) Close() error {
	return dhr.body.Close()
}

func (dhr *defaultHttpResponse) StatusCode() HttpStatus {
	return HttpStatus(dhr.statusCode)
}