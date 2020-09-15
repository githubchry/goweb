package drivers

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

var RedisDbConn redis.Conn
var RedisDbName string

// 初始化
func RedisDBInit() error {
	var err error
	RedisDbConn, err = redis.Dial("tcp", "172.20.209.220:6379")
	if err != nil {
		log.Fatal("Connect to redis error", err)
	}

	log.Println("Connected to RedisDB!")

	_, err = RedisDbConn.Do("AUTH", "chry")
	if err != nil {
		RedisDbConn.Close()
		log.Println("AUTH to RedisDB failed!")
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
