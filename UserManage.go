package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/githubchry/goweb/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Username	string 	`json:"username"`
	Password	string	`json:"password"`
	Email 		string	`json:"email"`
}


// AddReply represents the result of an addition operation.
type Reply struct {
	Result 	int		`json:"result"`
	Message string 	`json:"message"`
	Token 	string 	`json:"token"`
}

// [golang jwt-go的使用](https://www.cnblogs.com/jianga/p/12487267.html)
// [使用JWT进行接口认证](https://studygolang.com/articles/27242?fr=sidebar)
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	var user User
	// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
	json.NewDecoder(r.Body).Decode(&user)
	// 打印请求数据
	log.Println("post req: ", user.Username, user.Password, user.Email)

	// 操作mongodb
	// 返回结果
	var reply Reply

	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", user.Username).Decode(&result)

	if result.Password == user.Password {
		// 账号密码校验通过, 创建Token并返回
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))	//把当前时间秒数以十进制转成字符串写到h
		io.WriteString(h, user.Username)	// 除了写入时间, 再追加用户名
		token := fmt.Sprintf("%x", h.Sum(nil))
		log.Println("token", token)

		// 把token存到redis
		models.InsertToken(result.Username, token, 120);

		reply.Result = 0
		reply.Message = "登录成功!"
		reply.Token = token
	} else if len(result.Username) <= 0 {
		reply.Result = 1
		reply.Message = "用户名不存在!"
	} else {
		reply.Result = 2
		reply.Message = "密码错误!"
	}

	log.Println(reply.Message)
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(reply)

}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	var user User
	// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
	json.NewDecoder(r.Body).Decode(&user)
	// 打印请求数据
	log.Println("post req: ", user.Username, user.Password, user.Email)

	// 操作mongodb
	// 返回结果
	var reply Reply

	// 删除
	//deleteResult := models.NewMgo().DeleteMany("username", user.Username)
	//fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult)

	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", user.Username).Decode(&result)
	if len(result.Username) > 0 {
		// 用户名已注册
		reply.Result = 1
		reply.Message = "用户名已注册!"
	} else {
		// 单个插入
		InsertOneResult := models.NewMgo().InsertOne(user)
		log.Println("Inserted a single document: ", InsertOneResult)

		reply.Result = 0
		reply.Message = "注册成功!"
	}

	log.Println(reply.Message)
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(reply)
}