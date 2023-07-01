package fhttp

import "time"

type RequestWithRetryArgs struct {
	RequestParams      RequestParams
	RetryAttempts      uint
	RetryDelay         time.Duration
	ExpectedStatusCode int
}

type RequestParams struct {
	Method      string
	URL         string
	Proxy       string
	Body        []byte `json:"-"`
	QueryParams map[string]string
	Headers     map[string]string `json:"-"`
	Timeout     time.Duration
	Cookies     map[string]string
}
