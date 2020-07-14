package hammer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	urlpkg "net/url"
	"regexp"
	"strings"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

type UserAgentType int

const (
	Mobile UserAgentType = iota
	Pc
)

var UserAgentDefault = []string{
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
	"Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 6 Build/LYZ28E) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Mobile/14F89;GameHelper",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1",
	"Mozilla/5.0 (iPad; CPU OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:46.0) Gecko/20100101 Firefox/46.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.4",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:46.0) Gecko/20100101 Firefox/46.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/13.10586",
}

func UserAgent(ua ...UserAgentType) (res string) {
	var index int
	l := len(UserAgentDefault)
	if ua == nil || len(ua) == 0 {
		index = r.Intn(l - 1)
	} else if ua[0] == Mobile {
		index = r.Intn(6)
	} else if ua[0] == Pc {
		index = r.Intn(5) + 8
	} else {
		index = r.Intn(l - 1)
	}
	res = UserAgentDefault[index]
	return
}

func requestCloudMusicApi(method, url string, data map[string]interface{}, options *Options) (*http.Response, error) {
	method = strings.ToUpper(method)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	header := req.Header
	header.Add("User-Agent", UserAgent(options.Ua))
	if method == "POST" {
		header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if strings.Contains(url, "music.163.com") {
		header.Add("Referer", "https://music.163.com")
	}
	for _, v := range options.Cookies {
		req.AddCookie(v)
	}
	if header.Get("Cookie") == "" {
		header.Set("Cookie", options.Token)
	}
	if options.Crypto == "weapi" {
		var csrf_token string
		reg, _ := regexp.Compile(`_csrf=([^(;|$)]+)`)
		for _, v := range req.Cookies() {
			csrfCookie := reg.FindString(v.Name)
			if csrfCookie != "" {
				csrf_token = csrfCookie
				break
			}
		}
		data["csrf_token"] = csrf_token
		data = weapiEncrypt(data)
		reg, _ = regexp.Compile(`\w*api`)
		url = reg.ReplaceAllString(url, "weapi")
	} else if options.Crypto == "linuxapi" {
		reg, _ := regexp.Compile(`\w*api`)
		m := make(map[string]interface{})
		m["method"] = method
		m["url"] = reg.ReplaceAllString(url, "api")
		m["params"] = data
		data = linuxapiEncrypt(m)
		header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36")
		url = "https://music.163.com/api/linux/forward"
	} else if options.Crypto == "eapi" {
		var dataHeader = http.Header{}
		dataHeader.Add("osver", GetCookie(options.Cookies, "osver"))
		dataHeader.Add("deviceId", GetCookie(options.Cookies, "deviceId"))
		dataHeader.Add("appver", GetCookie(options.Cookies, "appver", "6.1.1"))
		dataHeader.Add("versioncode", GetCookie(options.Cookies, "versioncode", "140"))
		dataHeader.Add("mobilename", GetCookie(options.Cookies, "mobilename"))
		dataHeader.Add("buildver", GetCookie(options.Cookies, "buildver"))
		dataHeader.Add("resolution", GetCookie(options.Cookies, "resolution", "1920x1080"))
		dataHeader.Add("__csrf", GetCookie(options.Cookies, "__csrf"))
		dataHeader.Add("os", GetCookie(options.Cookies, "os", "android"))
		dataHeader.Add("channel", GetCookie(options.Cookies, "channel"))
		dataHeader.Add("channel", GetCookie(options.Cookies, "channel"))
		dataHeader.Add("requestId", fmt.Sprintf("%d_%04d", time.Now().UnixNano()/1000000, r.Intn(1000)))
		if c := GetCookie(options.Cookies, "MUSIC_U"); c != "" {
			dataHeader.Add("MUSIC_U", c)
		}
		if c := GetCookie(options.Cookies, "MUSIC_A"); c != "" {
			dataHeader.Add("MUSIC_A", c)
		}
		req.Header.Set("Cookie", "")
		for k, v := range dataHeader {
			req.AddCookie(&http.Cookie{
				Name:  k,
				Value: v[0],
			})
		}
		data["header"] = dataHeader
		data = eapiEncrypt(options.Url, data)
		reg, _ := regexp.Compile(`\w*api`)
		url = reg.ReplaceAllString(url, "eapi")
	}
	u, _ := urlpkg.Parse(url)
	req.URL = u
	req.Host = u.Host
	param := GetQueryParam(data)
	buf := new(bytes.Buffer)
	buf.Write([]byte(param))
	req.Body = ioutil.NopCloser(buf)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetCookie(cookies []*http.Cookie, name string, defaultValue ...string) string {
	for _, v := range cookies {
		if v.Name == name {
			return v.Value
		}
	}
	if defaultValue == nil || len(defaultValue) == 0 {
		return ""
	}
	return defaultValue[0]
}