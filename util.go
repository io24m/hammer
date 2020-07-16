package hammer

import (
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
