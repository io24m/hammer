package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Cfg struct {
	Port string
}

func (me *Cfg) Default() {
	me.Port = "8849"
}

func Config(name ...string) (cfg *Cfg) {
	cfg = &Cfg{}
	cfg.Default()
	path := "config.json"
	if name != nil && len(name) != 0 {
		path = name[0]
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = json.Unmarshal(file, cfg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
