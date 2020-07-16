package hammer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Cfg struct {
	Port            string
	SavePath        string
	PlayListId      string
	ConcurrentCount int
}

var defaultConfig = &Cfg{
	Port:            "8849",
	SavePath:        "\\",
	ConcurrentCount: 10,
}

func Config(name ...string) (cfg *Cfg) {
	cfg = defaultConfig
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
