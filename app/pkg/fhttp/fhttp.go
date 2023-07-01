package fhttp

import (
	"fmt"
	"server-template/config"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type Client struct {
	cfg *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (h *Client) Request(args RequestParams) (responseBody []byte, statusCode int, err error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(args.Method)
	for k, v := range args.Headers {
		req.Header.Set(k, v)
	}
	var queryCount int
	for k, v := range args.QueryParams {
		switch queryCount {
		case 0:
			args.URL += fmt.Sprintf("?%s=%s", k, v)
		default:
			args.URL += fmt.Sprintf("&%s=%s", k, v)
		}
		queryCount++
	}

	req.SetBody(args.Body)
	req.SetRequestURI(args.URL)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	client := &fasthttp.Client{}
	if args.Proxy != "" {
		client = &fasthttp.Client{
			Dial: fasthttpproxy.FasthttpHTTPDialer(args.Proxy),
		}
	}

	req.SetConnectionClose()

	if args.Timeout != 0 {
		err = client.DoTimeout(req, res, args.Timeout*time.Millisecond)
	} else {
		err = client.DoRedirects(req, res, 10)
	}
	if err != nil {
		return
	}

	statusCode = res.StatusCode()

	responseBody = make([]byte, len(res.Body()))
	copy(responseBody, res.Body())
	return
}
