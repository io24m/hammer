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
	get, err := http.Get(server.URL + "/login?email=xxx@163.com&password=yyy")
	if err != nil {
		panic(err)
	}
	defer get.Body.Close()
	t.Log(get.StatusCode)
	all, _ := ioutil.ReadAll(get.Body)
	t.Log(string(all))
}