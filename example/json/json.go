package main

import (
	jsonPkg "github.com/io24m/hammer/json"
	"log"
)

func main() {
	j := `{"a":1,"b":2,"c":[{"d":3},{"d":4}]}`
	json, err := jsonPkg.ReadJson(j)
	if err != nil {
		log.Fatal(err)
	}
	i, err := json.Get("a").Int()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(i == 1)
	i, err = json.Get("c[0].d").Int()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(i == 3)
	c := json.Get("c")
	i, err = c.Get("[0].d").Int()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(i == 3)
	i, err = c.Get("[1].d").Int()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(i == 4)
	_, err = c.Get("[2].d").Int()
	if err != nil {
		log.Fatal(err)
	}
}
