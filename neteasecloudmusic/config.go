package neteasecloudmusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type Cfg struct {
	Port            string
	SavePath        string
	PlayListId      string
	ConcurrentCount int
}

var defaultConfig *Cfg

func Config(name ...string) *Cfg {
	if defaultConfig != nil {
		return defaultConfig
	}
	r := &sync.RWMutex{}
	r.Lock()
	if defaultConfig != nil {
		return defaultConfig
	}
	defaultConfig = &Cfg{
		Port:            "8849",
		SavePath:        "\\",
		ConcurrentCount: 10,
	}
	r.Unlock()
	path := "config.json"
	if name != nil && len(name) != 0 {
		path = name[0]
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return defaultConfig
	}
	err = json.Unmarshal(file, defaultConfig)
	if err != nil {
		fmt.Println(err.Error())
		return defaultConfig
	}
	return defaultConfig
}
