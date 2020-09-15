package models

import (
	"github.com/githubchry/goweb/drivers"
	"github.com/gomodule/redigo/redis"
	"log"
)

// 插入token
func InsertToken(key string, val string, time int) error {

	var err error
	if time > 0 {
		_, err = drivers.RedisDbConn.Do("SET", key, val, "EX", time)
	} else {
		_, err = drivers.RedisDbConn.Do("SET", key, val)
	}

	if err != nil{
		log.Println("redis set", key, val, time," failed:", err)
	}

	return err
}

//
func FindToken(token string) (string, error) {

	return redis.String(drivers.RedisDbConn.Do("GET", token))
}
