package webapi

import (
	"encoding/json"
	"github.com/githubchry/goweb/models"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
)

type User struct {
	Username	string 	`json:"username"`
	Password	string	`json:"password"`
	Email 		string	`json:"email"`
}


// UserRsp represents the result of an addition operation.
type UserRsp struct {
	Result 	int		`json:"result"`
	Message string 	`json:"message"`
	Token 	string 	`json:"token"`
}

// [golang jwt-go的使用](https://www.cnblogs.com/jianga/p/12487267.html)
// [使用JWT进行接口认证](https://studygolang.com/articles/27242?fr=sidebar)
func UserLogin(w http.ResponseWriter, r *http.Request) {
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
	var rsp UserRsp

	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", user.Username).Decode(&result)

	if result.Password == user.Password {

		// 利用uuid库生成唯一且随机的token
		token := uuid.NewV4().String()
		log.Println("token", token)

		// 把token存到redis
		models.InsertToken(result.Username, token, 120);

		rsp.Result = 0
		rsp.Message = "登录成功!"
		rsp.Token = token
	} else if len(result.Username) <= 0 {
		rsp.Result = 1
		rsp.Message = "用户名不存在!"
	} else {
		rsp.Result = 2
		rsp.Message = "密码错误!"
	}

	log.Println(rsp.Message)
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(rsp)

}

func UserLogout(w http.ResponseWriter, r *http.Request) {

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

func UserRegister(w http.ResponseWriter, r *http.Request) {
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
	var rsp UserRsp

	// 删除
	//deleteResult := models.NewMgo().DeleteMany("username", user.Username)
	//fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult)

	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", user.Username).Decode(&result)
	if len(result.Username) > 0 {
		// 用户名已注册
		rsp.Result = 1
		rsp.Message = "用户名已注册!"
	} else {
		// 单个插入
		InsertOneResult := models.NewMgo().InsertOne(user)
		log.Println("Inserted a single document: ", InsertOneResult)

		rsp.Result = 0
		rsp.Message = "注册成功!"
	}

	log.Println(rsp.Message)
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(rsp)
}