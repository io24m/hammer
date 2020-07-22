package hammer

import (
	"net/http"
	"strings"
)

func Run() {
	cfg := Config()
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe("127.0.0.1:"+cfg.Port, nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	origin := r.Header.Get("origin")
	if strings.TrimSpace(origin) == "" {
		w.Header().Add("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Add("Access-Control-Allow-Origin", origin)
	}
	w.Header().Add("Access-Control-Allow-Headers", "X-Requested-With,Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(s))
}

var route map[string]func(*Query) (string, error)

func init() {
	route = make(map[string]func(*Query) (string, error))
	route["/login"] = Login
	route["/login/cellphone"] = LoginCellphone
	route["/playlist/detail"] = PlaylistDetail
	route["/song/detail"] = SongDetail
	route["/song/url"] = SongUrl
	route["/activate/init/profile"] = ActivateInitProfile
}
