package shared

import (
	"io"
	"net/http"
	"net/url"
)

type Query struct {
	Cookies []*http.Cookie
	Param   url.Values
	Body    io.ReadCloser
}

type Options struct {
	Crypto  string
	Cookies []*http.Cookie
	Proxy   interface{}
}

const (
	POST  string = "POST"
	LOGIN string = "https://music.163.com/weapi/login"
)
