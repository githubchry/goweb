package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//定义配置文件解析后的结构

type HTTPCfg struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

type GRPCCfg struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

type MinioCfg struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSL      bool   `json:"ssl"`
}

type MongoCfg struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type RedisCfg struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Network  string `json:"network"`
	Password string `json:"password"`
}

type KafkaCfg struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
}

type AppCfg struct {
	HTTPCfg  HTTPCfg  `json:"http"`
	GRPCCfg  GRPCCfg  `json:"grpc"`

	MinioCfg MinioCfg `json:"minio"`
	MongoCfg MongoCfg `json:"mongo"`
	RedisCfg RedisCfg `json:"redis"`
	KafkaCfg KafkaCfg `json:"kafka"`
}

/*************************************************
Function: LoadConfig
Description: read config file to config struct
@parameter filename: config file
Return: Config,bool
*************************************************/
func LoadConfig(filename string) (AppCfg, bool) {
	var conf AppCfg
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("ReadFile error")
		return conf, false
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Println("json.Unmarshal error")
		return conf, false
	}

	return conf, true
}

func SaveConfig(conf AppCfg, filename string) bool {

	data, err := json.Marshal(conf)
	if err != nil {
		log.Fatalln("json.Marshal error")
		return false
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatalln("WriteFile error")
		return false
	}

	return true
}
