package main

import (
	"encoding/json"
	"github.com/githubchry/goweb/models"
	"log"
	"net/http"
)

type FileReq struct {
	Cmd		int		`json:"cmd"`	// 0上传 1下载
	Type 	string 	`json:"type"`	// photo music
	Suffix	string 	`json:"suffix"`	// .jpg
}

// UserRsp represents the result of an addition operation.
type FileRsp struct {
	Result 	string	`json:"result"`
	Status 	int		`json:"status"`
	Message string 	`json:"message"`
}

func presignedUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	log.Println(r.Header.Get("Username"), r.Header.Get("Token"))

	var fileRsp FileRsp
	// 从redis查询token
	token, err := models.FindToken(r.Header.Get("Username"))
	if err != nil {
		log.Println("非法访问: Username/Token缺失或无效!")
		fileRsp.Status = -1
		fileRsp.Message = "非法访问: Username/Token缺失或无效!"
		json.NewEncoder(w).Encode(fileRsp)
		return
	}

	if token != r.Header.Get("Token") {
		log.Println("非法访问: Token无效或被重复登录顶号!")
		fileRsp.Status = -2
		fileRsp.Message = "非法访问: Token无效或被重复登录顶号!"
		json.NewEncoder(w).Encode(fileRsp)
		return
	}

	var fileReq FileReq

	// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
	json.NewDecoder(r.Body).Decode(&fileReq)

	filename := r.Header.Get("Username") + "." + fileReq.Suffix
	// 打印请求数据
	log.Println("post req: ", fileReq.Cmd, fileReq.Type, filename)

	if 0 == fileReq.Cmd {
		fileRsp.Result = models.PreUpload(fileReq.Type, filename)
	} else {
		fileRsp.Result = models.PreDownload(fileReq.Type, filename)
	}
	log.Println(fileRsp.Result)
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(fileRsp)

}

func download(w http.ResponseWriter, r *http.Request) {

	// 从redis查询token
	token, err := models.FindToken(r.Header.Get("Username"))
	if err != nil {
		return
	}

	if token != r.Header.Get("Token") {
		return
	}

	models.DeleteToken(r.Header.Get("Username"))
	log.Println("注销用户:", r.Header.Get("Username"))
	return
}
