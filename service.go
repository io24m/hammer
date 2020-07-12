package hammer

import (
	"io/ioutil"
)

func Login(query *Query) (string, error) {
	var options = &Options{
		Crypto:  "weapi",
		Cookies: query.Cookies,
		Proxy:   query.Proxy,
		Ua:      "pc",
	}
	var data interface{} = query.Param
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
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 {

	}
	return "", nil
}
