package hammer

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer server.Close()
	//get, err := http.Get(server.URL + "/login?email=xxx&password=yyy")
	//get, err := http.Get(server.URL + "/login/cellphone?phone=15831706253&password=yyy")
	get, err := http.Get(server.URL + "/playlist/detail?id=24381616")
	if err != nil {
		panic(err)
	}
	defer get.Body.Close()
	t.Log(get.StatusCode)
	all, _ := ioutil.ReadAll(get.Body)
	t.Log("all:==========")
	t.Log(string(all))
}
