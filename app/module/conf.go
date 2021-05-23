package module

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"sync"
)


type Config struct {
	Log *string
	Stream *string
	Region string
	Index string
	Esendpoint string
	Esuser string
	Espasswd string
}


var Configs Config
var once sync.Once
func init() {
	once.Do(func() {
		data,_ := ioutil.ReadFile("./conf/conf.yaml")
		err := yaml.Unmarshal(data,&Configs)
		if err !=nil {
			panic("decode error")
		}
	})
}