package hammer

import (
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

func DownPlayListSong() {
	cfg := Config()
	ids := getIds(cfg)
	fmt.Println("find:", len(ids))
	songs := getSongs(cfg, ids)
	names := getSongNames(cfg, ids)
	merSong(songs, names)
	mux := sync.Mutex{}
	completeCount := 0
	works := make(chan struct{}, cfg.ConcurrentCount)
	var wg sync.WaitGroup
	err := os.MkdirAll(cfg.SavePath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	for _, v := range songs {
		works <- struct{}{}
		wg.Add(1)
		go func(url, path string) {
			defer wg.Done()
			down(url, path)
			mux.Lock()
			completeCount++
			fmt.Println(strconv.Itoa(completeCount) + ":" + path)
			mux.Unlock()
			<-works
		}(v.url, cfg.SavePath+v.songName+`.`+v.songType)
	}
	wg.Wait()
	fmt.Println("complete")

}

type songDetails struct {
	songId   string
	songName string
	url      string
	songType string
}

func (me *songDetails) name(name string) {
	me.songName = name
}

func merSong(s1 map[string]*songDetails, s2 map[string]*songDetails) {
	for k, v := range s2 {
		details := s1[k]
		if details == nil {
			continue
		}
		details.songName = v.songName
	}
}

func getIds(cfg *Cfg) []string {
	query := &Query{}
	query.AddParam("id", cfg.PlayListId)
	resp, err := PlaylistDetail(query)
	if err != nil {
		panic(err)
	}
	trackIds := gojsonq.New().JSONString(resp).From("playlist.trackIds").Select("id").Get()
	return getIdArray(trackIds)
}

func getIdArray(i interface{}) (res []string) {
	ls := i.([]interface{})
	for _, v := range ls {
		m := v.(map[string]interface{})
		id := strconv.FormatFloat(m["id"].(float64), 'f', -1, 64)
		res = append(res, string(id))
	}
	return
}

func getSongs(cfg *Cfg, ids []string) (res map[string]*songDetails) {
	params := strings.Join(ids, ",")
	query := &Query{
		Param: url.Values{},
	}
	query.Param.Add("id", params)
	song, err := SongUrl(query)
	//respSongs, err := http.Get(cfg.host + "/song/url?id=" + params)
	if err != nil {
		panic(err)
	}
	datas := gojsonq.New().JSONString(song).From("data")
	res = id2url(datas.Select("id", "url", "type").Get())
	return
}

func getSongNames(cfg *Cfg, ids []string) (res map[string]*songDetails) {
	params := strings.Join(ids, ",")
	query := &Query{
		Param: url.Values{},
	}
	query.Param.Add("ids", params)
	detail, err := SongDetail(query)
	//respSongs, err := http.Get(cfg.host + "/song/detail?ids=" + params)
	if err != nil {
		panic(err)
	}
	datas := gojsonq.New().JSONString(detail).From("songs")
	idname := datas.Select("id", "name").Get()
	res = id2name(idname)
	return
}

func id2url(i interface{}) map[string]*songDetails {
	res := make(map[string]*songDetails)
	ls := i.([]interface{})
	for _, v := range ls {
		m := v.(map[string]interface{})
		id := strconv.FormatFloat(m["id"].(float64), 'f', -1, 64)
		if m["type"] == nil {
			continue
		}
		if m["url"] == nil {
			continue
		}

		res[id] = &songDetails{songId: id, url: m["url"].(string), songType: m["type"].(string)}
	}
	return res
}

func id2name(i interface{}) map[string]*songDetails {
	res := make(map[string]*songDetails)
	ls := i.([]interface{})
	for _, v := range ls {
		m := v.(map[string]interface{})
		id := strconv.FormatFloat(m["id"].(float64), 'f', -1, 64)
		res[id] = &songDetails{songId: id, songName: m["name"].(string)}
	}
	return res
}

func down(url, path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
}
