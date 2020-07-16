package hammer

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	UrlLogin          string = "https://music.163.com/weapi/login"
	UrlLoginCellphone string = "https://music.163.com/weapi/login/cellphone"
	UrlPlaylistDetail string = "https://music.163.com/weapi/v3/playlist/detail"
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
	if err == nil {
		return "", nil
	}
	return string(all), nil

}
