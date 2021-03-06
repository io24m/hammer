package neteasecloudmusic

import (
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

func init() {
	initRoute()
}

func Run() {
	cfg := Config()
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe("127.0.0.1:"+cfg.Port, nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
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
	f := route[r.URL.Path]
	if f == nil {
		w.Write([]byte("run success,eg:\n"))
		for _, v := range urls {
			w.Write([]byte(v))
		}
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
}

type requestFunc func(*Query) (string, error)

var route map[string]requestFunc
var urls []string

func funcNames(f []requestFunc) map[string]requestFunc {
	reg, _ := regexp.Compile(`[A-Z]`)
	res := make(map[string]requestFunc)
	for _, v := range f {
		name := funcName(v)
		name = reg.ReplaceAllStringFunc(name, func(s string) string {
			return "/" + strings.ToLower(s)
		})
		res[name] = v
	}
	return res
}

func funcName(f interface{}) string {
	fc := reflect.ValueOf(f)
	rF := runtime.FuncForPC(fc.Pointer())
	funcName := rF.Name()
	funcName = funcName[strings.LastIndex(funcName, ".")+1:]
	return funcName
}
