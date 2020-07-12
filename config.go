package hammer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

type Cfg struct {
	Port string
	Host string
}

var defaultConfig *Cfg

var lock sync.Mutex

func (me *Cfg) merCfg(other *Cfg) {
	if other == nil {
		return
	}
	if strings.TrimSpace(other.Port) != "" {
		me.Port = other.Port
	}
	if strings.TrimSpace(other.Host) != "" {
		me.Host = other.Host
	}
}

func Config(name ...string) (cfg *Cfg) {
	if defaultConfig != nil {
		return defaultConfig
	}
	lock.Lock()
	if defaultConfig != nil {
		return defaultConfig
	}
	defaultConfig = &Cfg{}
	lock.Unlock()
	cfg = defaultConfig
	cfg.Host = ""
	cfg.Port = "8849"
	path := "config.json"
	if name != nil && len(name) != 0 {
		path = name[0]
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	other := &Cfg{}
	err = json.Unmarshal(file, other)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cfg.merCfg(other)
	return
}
