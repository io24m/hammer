package hammer

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	urlAlbumDetailDynamic  string = "https://music.163.com/api/album/detail/dynamic"
	urlAlbumNewest         string = "https://music.163.com/api/discovery/newAlbum"
	urlAlbumSub            string = "https://music.163.com/api/album/%s"
	urlAlbumSublist        string = "https://music.163.com/weapi/album/sublist"
	urlArtistAlbum         string = "https://music.163.com/weapi/artist/albums/%s"
	urlArtistDesc          string = "https://music.163.com/weapi/artist/introduction"
)

//Login 邮箱登录
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
	query.AddCookie("os", "pc")
	var opt = &options{
		crypto:  weapi,
		ua:      pc,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	cmResult, err := requestCloudMusicApi(post, urlLogin, data, opt)
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

//LoginCellphone 电话登录
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
	query.AddCookie("os", "pc")
	opt := &options{
		crypto:  weapi,
		cookies: query.Cookies,
		proxy:   nil,
		ua:      pc,
	}
	cmResult, err := requestCloudMusicApi(post, urlLoginCellphone, data, opt)
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

//PlaylistDetail 歌单详情
func PlaylistDetail(query *Query) (string, error) {
	data := make(map[string]interface{}, 0)
	data["id"] = query.GetParam("id")
	data["n"] = 100000
	data["s"] = 8
	opt := &options{
		crypto:  linuxapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, urlPlaylistDetail, data, opt)
}

//SongUrl 歌曲链接
func SongUrl(query *Query) (string, error) {
	if MUSIC_U := query.GetCookie("MUSIC_U"); strings.TrimSpace(MUSIC_U) == "" {
		query.AddCookie("_ntes_nuid", hex.EncodeToString(key(16)))
	}
	query.AddCookie("os", "pc")
	data := make(map[string]interface{}, 0)
	data["ids"] = "[" + query.GetParam("id") + "]"
	if br := query.GetParam("br"); strings.TrimSpace(br) != "" {
		data["br"] = br
	}
	data["br"] = 999000
	opt := &options{
		crypto:  linuxapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, urlSongUrl, data, opt)
}

//SongDetail 歌曲详情
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
	opt := &options{
		crypto:  weapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, urlSongDetail, data, opt)
}

//ActivateInitProfile 初始化名字
func ActivateInitProfile(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["nickname"] = query.GetParam("nickname")
	opt := &options{
		crypto:  eapi,
		cookies: query.Cookies,
		url:     "/api/activate/initProfile",
	}
	return responseDefault(post, urlActivateInitProfile, data, opt)
}

//Album 专辑内容
func Album(query *Query) (string, error) {
	id := query.GetParam("id")
	opt := &options{
		crypto:  weapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, fmt.Sprintf(urlAlbum, id), nil, opt)
}

//AlbumDetailDynamic 专辑动态信息
func AlbumDetailDynamic(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["id"] = query.GetParam("id")
	opt := &options{
		crypto:  weapi,
		cookies: query.Cookies,
		proxy:   query.Proxy,
	}
	return responseDefault(post, urlAlbumDetailDynamic, data, opt)
}

//AlbumNewest 最新专辑
func AlbumNewest(query *Query) (string, error) {
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlAlbumNewest, nil, opt)
}

//AlbumSub 收藏/取消收藏专辑
func AlbumSub(query *Query) (string, error) {
	data := make(map[string]interface{})
	t := query.GetParam("t")
	data["id"] = query.GetParam("id")
	if t == "1" {
		t = "sub"
	} else {
		t = "unsub"
	}
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlAlbumSub, t), data, opt)
}

//AlbumSublist 已收藏专辑列表
func AlbumSublist(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["total"] = true
	if limit := query.GetParam("limit"); limit != "" {
		data["limit"] = limit
	} else {
		data["limit"] = 25
	}
	if offset := query.GetParam("offset"); offset != "" {
		data["offset"] = offset
	} else {
		data["offset"] = 0
	}
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlAlbumSublist, data, opt)
}

//ArtistAlbum 歌手专辑列表
func ArtistAlbum(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["total"] = true
	if limit := query.GetParam("limit"); limit != "" {
		data["limit"] = limit
	} else {
		data["limit"] = 30
	}
	if offset := query.GetParam("offset"); offset != "" {
		data["offset"] = offset
	} else {
		data["offset"] = 0
	}
	id := query.GetParam("id")
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, fmt.Sprintf(urlArtistAlbum, id), data, opt)
}

//ArtistDesc 歌手介绍
func ArtistDesc(query *Query) (string, error) {
	data := make(map[string]interface{})
	data["id"] = query.GetParam("id")
	opt := &options{crypto: weapi, cookies: query.Cookies, proxy: query.Proxy}
	return responseDefault(post, urlArtistDesc, data, opt)
}
