package models

import (
	"github.com/githubchry/goweb/internal/dao/drivers"
	"log"
	"net/url"
	"time"
)

// 获取上传url
func PreUpload(bucketName string, fileName string) string {

	presignedURL, err := drivers.MinioDbConn.PresignedPutObject(bucketName, fileName, time.Second * 24 * 60 * 60)
	if err != nil {
		log.Println(err)
	}
	//log.Println("Successfully generated presigned URL", presignedURL)
	return presignedURL.String()
}

// 获取下载url
func PreDownload(bucketName string, fileName string) string {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+fileName)

	presignedURL, err := drivers.MinioDbConn.PresignedGetObject(bucketName, fileName, time.Second * 60 * 2, reqParams)
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfully generated presigned URL", presignedURL)
	return presignedURL.String()
}

