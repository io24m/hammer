package hammer

import (
	"fmt"
	"github.com/io24m/hammer/jsongo"
	"github.com/io24m/hammer/mp3"
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
		go func(cfg *Cfg, song *songDetails) {
			defer wg.Done()
			//url := song.url
			path := cfg.SavePath + song.songName + `.` + song.songType
			down(cfg, song)
			mux.Lock()
			completeCount++
			fmt.Println(strconv.Itoa(completeCount) + ":" + path)
			mux.Unlock()
			<-works
		}(cfg, v)
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
	json, _ := jsongo.ReadJson(resp)
	ids := json.Get("playlist.trackIds").Map("id").Strings()
	return ids
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
	query := &Query{}
	query.AddParam("id", strings.Join(ids, ","))
	song, err := SongUrl(query)
	if err != nil {
		panic(err)
	}
	json, _ := jsongo.ReadJson(song)
	node := json.Get("data").Map().Values()
	res = id2url(node)
	return
}

func getSongNames(cfg *Cfg, ids []string) (res map[string]*songDetails) {
	params := strings.Join(ids, ",")
	query := &Query{
		Param: url.Values{},
	}
	query.Param.Add("ids", params)
	detail, err := SongDetail(query)
	if err != nil {
		panic(err)
	}
	json, _ := jsongo.ReadJson(detail)
	idname := json.Get("songs").Map().Values()
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

func down(cfg *Cfg, song *songDetails) {
	path := cfg.SavePath + song.songName + `.` + song.songType
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	resp, err := http.Get(song.url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	//_, err = io.Copy(f, resp.Body)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//test
	bytes, err := readBytes(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	v2_3, err := mp3.Mp3_ID3V2_3(bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	v2_3.Tag(mp3.TIT2, song.songName)
	fmt.Println(v2_3.Tags())
	_, err = f.Write(v2_3.Byte())
	if err != nil {
		fmt.Println(err)
		return
	}
}
