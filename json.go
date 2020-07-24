package hammer

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type IJson struct {
	data []byte
	json interface{}
}

func ReadJson(jsonStr string) (*IJson, error) {
	ijson := &IJson{
		data: []byte(jsonStr),
		json: make(map[string]interface{}),
	}
	err := json.Unmarshal(ijson.data, &ijson.json)
	if err != nil {
		return nil, err
	}
	return ijson, nil
}

func (json *IJson) Get(keys ...string) *IJsonNode {
	if len(keys) == 0 {
		return &IJsonNode{node: json.json}
	}
	key := keys[0]
	m, err := jsonKey(key)
	if err != nil {
		return &IJsonNode{}
	}
	return json.get(m)
}

func (json *IJson) get(m map[int]*nodeKey) *IJsonNode {
	var j interface{}
	j = json.json
	l := len(m)
	for i := 0; i < l; i++ {
		j = g(j, m[i])
	}
	return &IJsonNode{node: j}
}

func g(i interface{}, key *nodeKey) interface{} {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	v := m[key.name]
	if key.index == -1 {
		return v
	}
	s, ok := v.([]interface{})
	l := len(s)
	if !ok || l <= key.index {
		return v
	}
	//if l > 0 && l <= key.index {
	//	return s[0]
	//}
	return s[key.index]
}

type IJsonNodeList struct {
	nodes []*IJsonNode
}

func (jsonList *IJsonNodeList) Nodes() []*IJsonNode {
	return jsonList.nodes
}

func (jsonList *IJsonNodeList) Values() []interface{} {
	r := make([]interface{}, 0)
	for _, v := range jsonList.nodes {
		r = append(r, v.Value())
	}
	return r
}

func (jsonList *IJsonNodeList) Integers() ([]int64, error) {
	r := make([]int64, 0)
	for _, v := range jsonList.nodes {
		i, err := v.Int()
		if err != nil {
			return nil, nil
		}
		r = append(r, i)
	}
	return r, nil
}

func (jsonList *IJsonNodeList) Floats() ([]float64, error) {
	r := make([]float64, 0)
	for _, v := range jsonList.nodes {
		i, err := v.Float()
		if err != nil {
			return nil, nil
		}
		r = append(r, i)
	}
	return r, nil
}

func (jsonList *IJsonNodeList) Strings() []string {
	r := make([]string, 0)
	for _, v := range jsonList.nodes {
		r = append(r, v.String())
	}
	return r
}

func (jsonList *IJsonNodeList) Booleans() ([]bool, error) {
	r := make([]bool, 0)
	for _, v := range jsonList.nodes {
		i, err := v.Bool()
		if err != nil {
			return nil, nil
		}
		r = append(r, i)
	}
	return r, nil
}

type IJsonNode struct {
	node interface{}
}

func (node *IJsonNode) Map(keys ...string) *IJsonNodeList {
	list := &IJsonNodeList{}
	list.nodes = make([]*IJsonNode, 0)
	ls, ok := node.node.([]interface{})
	if !ok {
		list.nodes = append(list.nodes, node)
		return list
	}
	key := ""
	if len(keys) > 0 {
		key = keys[0]
	}
	for _, v := range ls {
		iJson := &IJson{json: v}
		jsonNode := iJson.Get(key)
		list.nodes = append(list.nodes, jsonNode)
	}
	return list
}

func (node *IJsonNode) Value() interface{} {
	return node.node
}

func (node *IJsonNode) Bool() (bool, error) {
	if b, ok := node.node.(bool); ok {
		return b, nil
	}
	return false, errors.New("not bool")
}

func (node *IJsonNode) Int() (int64, error) {
	if f, ok := node.node.(float64); ok {
		return int64(f), nil
	}
	return 0, errors.New("not int")
}

func (node *IJsonNode) Float() (float64, error) {
	if f, ok := node.node.(float64); ok {
		return f, nil
	}
	return 0, errors.New("not float")
}

func (node *IJsonNode) String() string {
	switch s := node.node.(type) {
	default:
		return ""
	case string:
		return s
	case float64:
		return strconv.FormatFloat(s, 'f', 0, 64)
	case bool:
		return strconv.FormatBool(s)
	}
}

type nodeKey struct {
	name  string
	index int
}

func jsonKey(key string) (map[int]*nodeKey, error) {
	m := make(map[int]*nodeKey)
	keys := strings.Split(key, ".")
	for i, v := range keys {
		if !strings.Contains(v, "[") {
			m[i] = &nodeKey{
				name:  v,
				index: -1,
			}
			continue
		}
		index0 := strings.Index(v, "[")
		index1 := strings.Index(v, "]")
		index, err := strconv.ParseInt(v[index0+1:index1], 10, 64)
		if err != nil {
			return nil, err
		}
		m[i] = &nodeKey{
			name:  v[0:index0],
			index: int(index),
		}
	}
	return m, nil
}
