package logics

import (
	"github.com/githubchry/goweb/internal/dao/models"
	"log"
	"path"
)

func PresignedUrl(req *FileReq) FileRsp {

	filename := req.Username + path.Ext(req.Filename)
	// 打印请求数据
	log.Println("post req: ", req.Cmd, req.Type, filename)
	var rsp FileRsp
	if 0 == req.Cmd {
		rsp.Url = models.PreUpload(req.Type, filename)
	} else {
		rsp.Url = models.PreDownload(req.Type, filename)
	}
	log.Println(rsp.Url)

	return rsp
}