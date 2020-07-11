package hammer

import (
	"fmt"
	"github.com/io24m/hammer/service"
	"github.com/io24m/hammer/shared"
	"github.com/io24m/hammer/util"
	"net/http"
)

func Run() {
	cfg := util.Config()
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe("127.0.0.1:"+cfg.Port, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr, "连接成功")
	// 请求方式：GET POST DELETE PUT UPDATE
	fmt.Println("method:", r.Method)
	// /go
	fmt.Println("url:", r.URL.Path)
	fmt.Println("url:", r.URL.Query())
	fmt.Println("header:", r.Header)
	fmt.Println("body:", r.Body)
	//login := service.Login(nil)
	f := url2func[r.URL.Path]
	if f == nil {
		w.Write([]byte("启动服务成功"))
		return
	}
	r.Cookies()
	var query = &shared.Query{
		Cookies: r.Cookies(),
		Body:    r.Body,
		Param:   r.URL.Query(),
	}
	s, err := f(query)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(s))
	return
}

var url2func map[string]func(*shared.Query) (string, error)

func init() {
	url2func = make(map[string]func(*shared.Query) (string, error))
	url2func["/login"] = service.Login
}
