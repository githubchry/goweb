package drivers

import (
	"github.com/githubchry/goweb/configs"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
)

var RedisDbConn redis.Conn
var RedisDbName string

// 初始化
func RedisDBInit(cfg configs.RedisCfg) error {
	var err error
	RedisDbConn, err = redis.Dial(cfg.Network, cfg.Addr+":"+strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatal("Connect to redis error", err)
	}

	log.Println("Connected to RedisDB!")

	_, err = RedisDbConn.Do("AUTH", cfg.Password)
	if err != nil {
		log.Fatal("AUTH to RedisDB failed!")
	}

	return err
}

// 关闭
func RedisDBExit() {
	err := RedisDbConn.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("RedisDB is closed.")
}
