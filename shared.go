package hammer

import (
	"io"
	"net/http"
	"net/url"
)

type Query struct {
	Cookies []*http.Cookie
	Param   url.Values
	Body    io.ReadCloser
	Proxy   interface{}
}

type Options struct {
	Crypto  string
	Cookies []*http.Cookie
	Proxy   interface{}
	Ua      UserAgentType
	Token   string
	Url     string
}
