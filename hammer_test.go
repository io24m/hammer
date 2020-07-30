package hammer

import (
	"fmt"
	mp3pkg "github.com/io24m/hammer/mp3"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	initRoute()
	server := httptest.NewServer(http.HandlerFunc(indexHandler))
	defer server.Close()
	//get, err := http.Get(server.URL + "/login?email=xxx&password=yyy")
	//get, err := http.Get(server.URL + "/login/cellphone?phone=15831706220&md5_password=c42c9549055bcae217fecdb249fbc6a8")
	//get, err := http.Get(server.URL + "/playlist/detail?id=24381616")
	//get, err := http.Get(server.URL + "/song/detail?ids=347230,15")
	//get, err := http.Get(server.URL + "/activate/init/profile?nickname=testUser2019")
	//get, err := http.Get(server.URL + "/album/newest")
	get, err := http.Get(server.URL + "/artist/desc?id=5770")
	if err != nil {
		panic(err)
	}
	defer get.Body.Close()
	t.Log(get.StatusCode)
	all, _ := ioutil.ReadAll(get.Body)
	t.Log(string(all))
}

func TestReadJson(t *testing.T) {
	var s interface{}
	j, _ := ReadJson(`{"q":"w","ss":[{"ss":null,"a":1,"b":2,"C":[1,2],"d":[{"e":"qw"},{"e":"qw"}]}]}`)
	s = j.Get("q").String()
	fmt.Println(s)
	s = j.Get("w").String()
	fmt.Println(s)
	s = j.Get("C[1]").String()
	fmt.Println(s)
	s = j.Get("d[1].e").String()
	fmt.Println(s)
	s = j.Get("ss").Map("d").Nodes()[0].Map("e").Values()
	fmt.Println(s)
}

func TestMp3Read(t *testing.T) {
	test, _ := os.Open(`C:\down\NeteaseCloudMusic\test.mp3`)
	defer test.Close()
	testTitle, _ := os.Create(`C:\down\NeteaseCloudMusic\testTitle.mp3`)
	defer testTitle.Close()
	bytes, _ := readBytes(test)
	mp3, _ := mp3pkg.Mp3_ID3V2_3(bytes)
	//mp3.Tag(mp3pkg.TIT2, "test")
	testTitle.Write(mp3.Byte())
}

func TestDownPlayListSong(t *testing.T) {
	cfg := Config()
	cfg.SavePath = `\down\NeteaseCloudMusic\temp\`
	cfg.ConcurrentCount = 1
	DownPlayListSong()
}
