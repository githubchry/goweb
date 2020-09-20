package drivers

import (
	"log"
)

func init()  {
	// 初始化连接到MongoDB
	err := MongoDBInit()
	if err != nil {
		log.Fatal(err)
		return
	}

	// 初始化连接到RedisDB
	err = RedisDBInit()
	if err != nil {
		log.Fatal(err)
		return
	}

	// 初始化连接到MinioDB
	err = MinioDBInit()
	if err != nil {
		log.Fatal(err)
		return
	}
}
