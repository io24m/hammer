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
	f := route[r.URL.Path]
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

var route map[string]func(*Query) (string, error)

func init() {
	route = make(map[string]func(*Query) (string, error))
	route["/login"] = Login
	route["/login/cellphone"] = LoginCellphone
	route["/playlist/detail"] = PlaylistDetail
}
