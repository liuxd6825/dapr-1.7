package http

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func testServer() {
	url := "http://localhost:3500/v1.0/event-storage/aggregates/001/04-00050"
	c := newClient()
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second * 3)
		c.httpGet(context.Background(), url).OnError(func(err error) {
			println(fmt.Sprintf("test http server %d error :%s", i, err))
		}).OnSuccess(func(data []byte) {
			println(fmt.Sprintf("test http service %d success data:%s", i, data))
		})
	}
}

type client struct {
	httpClient *http.Client
}

func newClient() *client {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     time.Second * time.Duration(10),
		},
	}
	return &client{
		httpClient: httpClient,
	}
}

func (c *client) httpGet(ctx context.Context, url string) *Response {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return NewResponse(nil, err)
	}
	bs, err := c.getBodyBytes(resp)
	if resp.StatusCode != http.StatusOK {
		return NewResponse(nil, err)
	}
	return NewResponse(bs, err)
}

func (c *client) getBodyBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	return bs, err
}

type Response struct {
	data []byte
	err  error
}

func NewResponse(data []byte, err error) *Response {
	return &Response{
		data: data,
		err:  err,
	}
}

func (r *Response) GetData() []byte {
	return r.data
}

func (r *Response) GetError() error {
	return r.err
}

func (r *Response) OnSuccess(fun func(data []byte)) *Response {
	if r.data == nil {
		return r
	}
	if r.err != nil {
		return r
	}
	fun(r.data)
	return r
}

func (r *Response) OnError(fun func(err error)) *Response {
	if r.err != nil {
		fun(r.err)
	}
	return r
}
