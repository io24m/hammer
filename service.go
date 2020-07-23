package hammer

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	post                   string = "POST"
	get                    string = "GET"
	urlLogin               string = "https://music.163.com/weapi/login"
	urlLoginCellphone      string = "https://music.163.com/weapi/login/cellphone"
	urlPlaylistDetail      string = "https://music.163.com/weapi/v3/playlist/detail"
	urlSongDetail          string = "https://music.163.com/weapi/v3/song/detail"
	urlSongUrl             string = "https://music.163.com/api/song/enhance/player/url"
	urlActivateInitProfile string = "http://music.163.com/eapi/activate/initProfile"
	urlAlbum               string = "https://music.163.com/weapi/v1/album/%s"
)

func Login(query *Query) (string, error) {
	var data = make(map[string]interface{})
	data["username"] = query.GetParam("email")
	data["rememberLogin"] = "true"
	if md5Password := query.GetParam("md5_password"); strings.TrimSpace(md5Password) != "" {
		data["password"] = md5Password
	} else {
		pw := query.GetParam("password")
		sum := md5.Sum([]byte(pw))
		data["password"] = hex.EncodeToString(sum[:])
	}

	query.Cookies = append(query.Cookies, &http.Cookie{
		Name:  "os",
		Value: "pc",
	})
	var options = &options{
		crypto:  weapi,
		ua:      pc,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	cmResult, err := requestCloudMusicApi(post, urlLogin, data, options)
	if err != nil {
		return "", err
	}
	defer cmResult.Body.Close()
	body, err := ioutil.ReadAll(cmResult.Body)
	if err != nil {
		return "", err
	}
	m := make(map[string]interface{}, 0)
	json.Unmarshal(body, &m)
	float := strconv.FormatFloat(m["code"].(float64), 'f', -1, 64)
	if float == "502" {
		return "账号或密码错误", nil
	}
	m["cookie"] = cmResult.Header.Get("Set-Cookie")
	marshal, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	//code502	var msg = `{'msg':'账号或密码错误','code':'502','message':'账号或密码错误'}`
	return string(marshal), nil
}

func LoginCellphone(query *Query) (string, error) {
	data := make(map[string]interface{}, 0)
	data["phone"] = query.GetParam("phone")
	if cc := query.GetParam("countrycode"); strings.TrimSpace(cc) != "" {
		data["countrycode"] = query.GetParam("countrycode")
	}
	data["rememberLogin"] = "true"
	if md5Password := query.GetParam("md5_password"); strings.TrimSpace(md5Password) != "" {
		data["password"] = md5Password
	} else {
		pw := query.GetParam("password")
		sum := md5.Sum([]byte(pw))
		data["password"] = hex.EncodeToString(sum[:])
	}
	query.Cookies = append(query.Cookies, &http.Cookie{
		Name:  "os",
		Value: "pc",
	})
	options := &options{
		crypto:  weapi,
		cookies: query.Cookies,
		proxy:   nil,
		ua:      pc,
	}
	cmResult, err := requestCloudMusicApi(post, urlLoginCellphone, data, options)
	if err != nil {
		return "", err
	}
	defer cmResult.Body.Close()
	body, err := ioutil.ReadAll(cmResult.Body)
	if err != nil {
		return "", err
	}
	m := make(map[string]interface{}, 0)
	json.Unmarshal(body, &m)
	m["cookie"] = cmResult.Header.Get("Set-Cookie")
	marshal, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}

func PlaylistDetail(query *Query) (string, error) {
	data := make(map[string]interface{}, 0)
	data["id"] = query.GetParam("id")
	data["n"] = 100000
	data["s"] = 8
	options := &options{
		crypto:  linuxapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, urlPlaylistDetail, data, options)
}

func SongUrl(query *Query) (string, error) {
	if MUSIC_U := getCookie(query.Cookies, "MUSIC_U"); strings.TrimSpace(MUSIC_U) == "" {
		query.Cookies = addCookie(query.Cookies, "_ntes_nuid", hex.EncodeToString(key(16)))
	}
	query.Cookies = addCookie(query.Cookies, "os", "pc")
	data := make(map[string]interface{}, 0)
	data["ids"] = "[" + query.GetParam("id") + "]"
	if br := query.GetParam("br"); strings.TrimSpace(br) != "" {
		data["br"] = br
	}
	data["br"] = 999000
	options := &options{
		crypto:  linuxapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, urlSongUrl, data, options)
}

func SongDetail(query *Query) (string, error) {
	ids := query.GetParam("ids")
	reg, _ := regexp.Compile(`\s*,\s*`)
	idList := reg.Split(ids, -1)
	c := make([]string, 0)
	for _, v := range idList {
		c = append(c, fmt.Sprintf(`{"id":%s}`, v))
	}
	data := make(map[string]interface{})
	data["c"] = "[" + strings.Join(c, ",") + "]"
	data["ids"] = "[" + strings.Join(idList, ",") + "]"
	options := &options{
		crypto:  weapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, urlSongDetail, data, options)
}

func ActivateInitProfile(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["nickname"] = query.GetParam("nickname")
	options := &options{
		crypto:  eapi,
		cookies: query.Cookies,
		url:     "/api/activate/initProfile",
	}
	return responseDefault(post, urlActivateInitProfile, data, options)
}

func Album(query *Query) (string, error) {
	id := query.GetParam("id")
	options := &options{
		crypto:  weapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, fmt.Sprintf(urlAlbum, id), nil, options)
}
