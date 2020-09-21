package drivers

import (
	"github.com/githubchry/goweb/configs"
	"github.com/minio/minio-go"
	"log"
	"strconv"
)

var MinioDbConn *minio.Client
var MinioDbName string

// 初始化
func MinioDBInit(cfg configs.MinioCfg) error {
	var err error
	// Minio client需要以下4个参数来连接与Amazon S3兼容的对象存储。
	endpoint := cfg.Addr + ":" + strconv.Itoa(cfg.Port) // 对象存储服务的URL
	accessKeyID := cfg.Username                         //Access key是唯一标识你的账户的用户ID。
	secretAccessKey := cfg.Password                     //Secret key是你账户的密码。
	useSSL := cfg.SSL                                   //true代表使用HTTPS

	// 初使化 minio client对象。
	MinioDbConn, err = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to MinioDB!")

	location := "us-east-1"
	bucketNameArr := [...]string{"music", "photo"}

	//range遍历数组
	for _, bucketName := range bucketNameArr {
		// 创建存储桶。
		err = MinioDbConn.MakeBucket(bucketName, location)
		if err != nil {
			// 检查存储桶是否已经存在。
			var exists bool
			exists, err = MinioDbConn.BucketExists(bucketName)
			if err == nil && exists {
				log.Printf("We already own %s\n", bucketName)
			} else {
				log.Fatalln(err)
			}
		}
		log.Printf("Successfully created %s\n", bucketName)
	}
	return err
}

// 关闭
func MinioDBExit() {
	MinioDbConn = nil
	log.Println("MinioDB is closed.")
}
