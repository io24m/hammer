package shared

type Query struct {
	Cookie interface{}
	Param  interface{}
	Body   interface{}
}

type ServiceUrl string

const (
	POST  string     = "POST"
	LOGIN ServiceUrl = "https://music.163.com/weapi/login"
)
