package hammer

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

func queryParam(data interface{}) (res []string) {
	res = make([]string, 0)
	m := data.(map[string]interface{})
	if m == nil {
		return res
	}
	for k, v := range m {
		if vm, ok := v.(map[string]interface{}); ok {
			param := queryParam(vm)
			res = append(res, param...)
			break
		}
		if vs, ok := v.(string); ok {
			res = append(res, k+"="+url.QueryEscape(vs))
		}
	}
	return
}

func queryParamString(data interface{}) string {
	param := queryParam(data)
	join := strings.Join(param, "&")
	return join
}

func getCookie(cookies []*http.Cookie, name string, defaultValue ...string) string {
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

func readBytes(reader io.Reader) ([]byte, error) {
	var chunk []byte
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		chunk = append(chunk, buf[:n]...)
	}
	return chunk, nil
}
