package webapi

import (
	"encoding/json"
	"github.com/githubchry/goweb/models"
	"log"
	"net/http"
	"strconv"
)

// AddReq represents the parameter of an addition operation.
type AddReq struct {
	OperandA int
	OperandB int
}

// AddRsp represents the result of an addition operation.
type AddRsp struct {
	Result int		`json:"result"`	// 结果 OperandA + OperandB
	Status int		`json:"status"`	// 状态 0表示sucess
	Message string	`json:"message"`// 消息
}

func Addpost(w http.ResponseWriter, r *http.Request) {

	//log.Println(r.Header.Get("Username"), r.Header.Get("Token"))

	var addRsp AddRsp
	// 从redis查询token
	token, err := models.FindToken(r.Header.Get("Username"))
	if err != nil {
		log.Println("非法访问: Username/Token缺失或无效!")
		addRsp.Status = -1
		addRsp.Message = "非法访问: Username/Token缺失或无效!"
		json.NewEncoder(w).Encode(addRsp)
		return
	}

	if token != r.Header.Get("Token") {
		log.Println("非法访问: Token无效或被重复登录顶号!")
		addRsp.Status = -2
		addRsp.Message = "非法访问: Token无效或被重复登录顶号!"
		json.NewEncoder(w).Encode(addRsp)
		return
	}

	// 延长token生存周期?

	var addReq AddReq
	// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
	json.NewDecoder(r.Body).Decode(&addReq)
	// 打印请求数据
	log.Println("post req: ", addReq.OperandA, addReq.OperandB)
	// 进行加法计算，并保存结果到结构体内
	addRsp.Result = addReq.OperandA + addReq.OperandB
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(addRsp)
}

func Addget(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Header.Get("Username"), r.Header.Get("Token"))

	var addRsp AddRsp
	// 从redis查询token
	token, err := models.FindToken(r.Header.Get("Username"))
	if err != nil {
		log.Println("非法访问: Username/Token缺失或无效!")
		addRsp.Status = -1
		addRsp.Message = "非法访问: Username/Token缺失或无效!"
		json.NewEncoder(w).Encode(addRsp)
		return
	}

	if token != r.Header.Get("Token") {
		log.Println("非法访问: Token无效或被重复登录顶号!")
		addRsp.Status = -2
		addRsp.Message = "非法访问: Token无效或被重复登录顶号!"
		json.NewEncoder(w).Encode(addRsp)
		return
	}

	values := r.URL.Query()

	var addReq AddReq
	addReq.OperandA, _ = strconv.Atoi(values.Get("OperandA"))
	addReq.OperandB, _ = strconv.Atoi(values.Get("OperandB"))

	// 打印请求数据
	log.Println("get req: ", addReq.OperandA, addReq.OperandB)

	// 进行加法计算，并保存结果到结构体内
	addRsp.Result = addReq.OperandA + addReq.OperandB
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(addRsp)
}