package logics

import (
	"github.com/githubchry/goweb/internal/dao/models"
	"log"
)

type FileReq struct {
	Cmd		int		`json:"cmd"`	// 0上传 1下载
	Type 	string 	`json:"type"`	// photo music
	Suffix	string 	`json:"suffix"`	// .jpg
	Username string `json:"username"`	// .jpg
}

// UserRsp represents the result of an addition operation.
type FileRsp struct {
	Result 	string	`json:"result"`
	Status 	Status		`json:"status"`
	Message string 	`json:"message"`
}

func PresignedUrl(fileReq FileReq) FileRsp {

	filename := fileReq.Username + "." + fileReq.Suffix
	// 打印请求数据
	log.Println("post req: ", fileReq.Cmd, fileReq.Type, filename)
	var fileRsp FileRsp
	if 0 == fileReq.Cmd {
		fileRsp.Result = models.PreUpload(fileReq.Type, filename)
	} else {
		fileRsp.Result = models.PreDownload(fileReq.Type, filename)
	}
	log.Println(fileRsp.Result)

	return fileRsp
}