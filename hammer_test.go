package hammer

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func t1() {

}

func TestRun(t *testing.T) {
	initRoute()
	server := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer server.Close()
	//get, err := http.Get(server.URL + "/login?email=xxx&password=yyy")
	//get, err := http.Get(server.URL + "/login/cellphone?phone=15831706220&md5_password=c42c9549055bcae217fecdb249fbc6a8")
	//get, err := http.Get(server.URL + "/playlist/detail?id=24381616")
	get, err := http.Get(server.URL + "/song/detail?ids=347230,15")
	//get, err := http.Get(server.URL + "/activate/init/profile?nickname=testUser2019")
	if err != nil {
		panic(err)
	}
	defer get.Body.Close()
	t.Log(get.StatusCode)
	all, _ := ioutil.ReadAll(get.Body)
	t.Log(string(all))
}
