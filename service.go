package hammer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Login(query *Query) (string, error) {

	var data = make(map[string]interface{})
	data["username"] = query.Param.Get("email")
	if md5Password := query.Param.Get("md5_password"); strings.TrimSpace(md5Password) != "" {
		data["password"] = md5Password
	} else {
		pw := query.Param.Get("password")
		hash := md5.New()
		hash.Write([]byte(pw))
		data["password"] = hex.EncodeToString(hash.Sum(nil))
	}
	query.Cookies = append(query.Cookies, &http.Cookie{
		Name:  "os",
		Value: "pc",
	})
	var options = &Options{
		Crypto:  "weapi",
		Cookies: query.Cookies,
		Proxy:   query.Proxy,
		Ua:      Pc,
	}
	request := CreatRequest(POST, LOGIN, data, options)
	resp, err := request()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 502 {
		var msg = `{'msg':'账号或密码错误','code':'502','message':'账号或密码错误'}`
		return msg, nil
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 {

	}
	fmt.Println("login 请求成功：" + string(res))
	return string(res), nil
}
