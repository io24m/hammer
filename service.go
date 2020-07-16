package hammer

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	UrlLogin          string = "https://music.163.com/weapi/login"
	UrlLoginCellphone string = "https://music.163.com/weapi/login/cellphone"
	UrlPlaylistDetail string = "https://music.163.com/weapi/v3/playlist/detail"
	UrlSongDetail     string = "https://music.163.com/weapi/v3/song/detail"
	UrlSongUrl        string = "https://music.163.com/api/song/enhance/player/url"
)

func Login(query *Query) (string, error) {
	var data = make(map[string]interface{})
	data["username"] = query.Param.Get("email")
	data["rememberLogin"] = "true"
	if md5Password := query.Param.Get("md5_password"); strings.TrimSpace(md5Password) != "" {
		data["password"] = md5Password
	} else {
		pw := query.Param.Get("password")
		sum := md5.Sum([]byte(pw))
		data["password"] = hex.EncodeToString(sum[:])
	}

	query.Cookies = append(query.Cookies, &http.Cookie{
		Name:  "os",
		Value: "pc",
	})
	var options = &Options{
		Crypto:  "weapi",
		Ua:      Pc,
		Cookies: query.Cookies,
		Proxy:   query.Proxy,
	}
	cmResult, err := requestCloudMusicApi(POST, UrlLogin, data, options)
	if err != nil {
		return "", err
	}
	defer cmResult.Body.Close()
	body, err := ioutil.ReadAll(cmResult.Body)
	if err != nil {
		return "", err
	}
	//code502	var msg = `{'msg':'账号或密码错误','code':'502','message':'账号或密码错误'}`
	return string(body), nil
}

func LoginCellphone(query *Query) (string, error) {
	data := make(map[string]interface{}, 0)
	data["phone"] = query.Param.Get("phone")
	if cc := query.Param.Get("countrycode"); strings.TrimSpace(cc) != "" {
		data["countrycode"] = query.Param.Get("countrycode")
	}
	data["rememberLogin"] = "true"
	if md5Password := query.Param.Get("md5_password"); strings.TrimSpace(md5Password) != "" {
		data["password"] = md5Password
	} else {
		pw := query.Param.Get("password")
		sum := md5.Sum([]byte(pw))
		data["password"] = hex.EncodeToString(sum[:])
	}
	query.Cookies = append(query.Cookies, &http.Cookie{
		Name:  "os",
		Value: "pc",
	})
	options := &Options{
		Crypto:  "weapi",
		Cookies: query.Cookies,
		Proxy:   nil,
		Ua:      Pc,
	}
	cmResult, err := requestCloudMusicApi(POST, UrlLoginCellphone, data, options)
	if err != nil {
		return "", err
	}
	defer cmResult.Body.Close()
	body, err := ioutil.ReadAll(cmResult.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func PlaylistDetail(query *Query) (string, error) {
	data := make(map[string]interface{}, 0)
	data["id"] = query.Param.Get("id")
	data["n"] = 100000
	data["s"] = 8
	options := &Options{
		Crypto:  "linuxapi",
		Cookies: query.Cookies,
		Proxy:   query.Proxy,
	}
	api, err := requestCloudMusicApi(POST, UrlPlaylistDetail, data, options)
	if err != nil {
		return "", nil
	}
	defer api.Body.Close()
	all, err := ioutil.ReadAll(api.Body)
	if err != nil {
		return "", nil
	}
	m := make(map[string]interface{}, 0)
	json.Unmarshal(all, &m)
	return string(all), nil

}
func SongUrl(query *Query) (string, error) {
	if MUSIC_U := getCookie(query.Cookies, "MUSIC_U"); strings.TrimSpace(MUSIC_U) == "" {
		query.Cookies = addCookie(query.Cookies, "_ntes_nuid", hex.EncodeToString(key(16)))
	}
	query.Cookies = addCookie(query.Cookies, "os", "pc")
	data := make(map[string]interface{}, 0)
	data["ids"] = "[" + query.Param.Get("id") + "]"
	if br := query.Param.Get("br"); strings.TrimSpace(br) != "" {
		data["br"] = br
	}
	data["br"] = 999000
	options := &Options{
		Crypto:  "linuxapi",
		Cookies: query.Cookies,
		Proxy:   query.Proxy,
	}
	res, err := requestCloudMusicApi(POST, UrlSongUrl, data, options)
	if err != nil {
		return "", nil
	}
	defer res.Body.Close()
	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}
	return string(all), nil
}

func SongDetail(query *Query) (string, error) {
	ids := query.Param.Get("ids")
	reg, _ := regexp.Compile(`\s*,\s*`)
	idList := reg.Split(ids, -1)
	c := make([]string, 0)
	for _, v := range idList {
		c = append(c, fmt.Sprintf(`{"id":%s}`, v))
	}
	data := make(map[string]interface{}, 0)
	data["c"] = "[" + strings.Join(c, ",") + "]"
	data["ids"] = "[" + strings.Join(idList, ",") + "]"
	options := &Options{
		Crypto:  "weapi",
		Cookies: query.Cookies,
		Proxy:   query.Proxy,
		Ua:      0,
		Token:   "",
		Url:     "",
	}
	res, err := requestCloudMusicApi(POST, UrlSongDetail, data, options)
	if err != nil {
		return "", nil
	}
	defer res.Body.Close()
	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}

	return string(all), nil

}
