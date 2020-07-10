package hammer

import (
	"fmt"
	"github.com/io24m/hammer/service"
	"net/http"
)

func Run() {
	http.HandleFunc("/go", myHandler)
	http.ListenAndServe("127.0.0.1:8849", nil)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr, "连接成功")
	// 请求方式：GET POST DELETE PUT UPDATE
	fmt.Println("method:", r.Method)
	// /go
	fmt.Println("url:", r.URL.Path)
	fmt.Println("header:", r.Header)
	fmt.Println("body:", r.Body)
	login := service.Login(nil)
	// 回复
	w.Write([]byte(login))
}
