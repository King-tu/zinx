package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	IP string
	Port uint32
	Name string
	IPVerson string
	WorkerSize int
	TaskQueneSize int
	MaxConnNum int
}

var GlobalConfig *Config

func (conf *Config)Reload()  {
	confInfo, err := ioutil.ReadFile("./conf/conf.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(confInfo, &GlobalConfig)
	if err != nil {
		log.Fatalln(err)
	}
}

func init()  {

	//默认配置
	GlobalConfig := &Config{
		IP:       "127.0.0.1",
		Port:     8080,
		Name:     "zinx Server",
		IPVerson: "tcp4",
		WorkerSize: 0,
		TaskQueneSize: 0,
		MaxConnNum: 1,
	}

	//重新加载用户配置
	GlobalConfig.Reload()
}
