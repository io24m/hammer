package hammer

import (
	"net/http"
)

func Run() {
	cfg := Config()
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe("127.0.0.1:"+cfg.Port, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.RemoteAddr, "连接成功")
	f := url2func[r.URL.Path]
	if f == nil {
		w.Write([]byte("启动服务成功"))
		return
	}

	s, err := f(&Query{
		Cookies: r.Cookies(),
		Body:    r.Body,
		Param:   r.URL.Query(),
	})
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(s))
	return
}

var url2func map[string]func(*Query) (string, error)

func init() {
	url2func = make(map[string]func(*Query) (string, error))
	url2func["/login"] = Login
}
