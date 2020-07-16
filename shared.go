package hammer

import (
	"io"
	"net/http"
	"net/url"
)

const (
	mobile userAgentType = "mobile"
	pc     userAgentType = "pc"
)

type userAgentType string

type Query struct {
	Cookies []*http.Cookie
	Param   url.Values
	Body    io.ReadCloser
	Proxy   interface{}
}

func (query *Query) AddParam(name, value string) *Query {
	if query.Param == nil {
		query.Param = url.Values{}
	}
	query.Param.Add(name, value)
	return query
}

type Options struct {
	Crypto  string
	Cookies []*http.Cookie
	Proxy   interface{}
	Ua      userAgentType
	Token   string
	Url     string
}
