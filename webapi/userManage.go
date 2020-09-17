package webapi

import (
	"encoding/json"
	"fmt"
	"github.com/githubchry/goweb/models"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	Username	string 	`json:"username"`
	Password	string	`json:"password"`
	Email 		string	`json:"email"`
	Photo 		string	`json:"photo"`
}


// UserRsp represents the result of an addition operation.
type UserRsp struct {
	Status 	int		`json:"status"`			//应答状态
	Result 	string	`json:"result"`
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

		rsp.Status = 0
		rsp.Message = "登录成功!"
		rsp.Token = token
	} else if len(result.Username) <= 0 {
		rsp.Status = 1
		rsp.Message = "用户名不存在!"
	} else {
		rsp.Status = 2
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
		rsp.Status = 1
		rsp.Message = "用户名已注册!"
	} else {
		// 单个插入
		InsertOneResult := models.NewMgo().InsertOne(user)
		log.Println("Inserted a single document: ", InsertOneResult)

		rsp.Status = 0
		rsp.Message = "注册成功!"
	}

	log.Println(rsp.Message)
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(rsp)
}

func UserSetPhoto(w http.ResponseWriter, r *http.Request) {

	var userRsp UserRsp
	// 从redis查询token
	token, err := models.FindToken(r.Header.Get("Username"))
	if err != nil {
		log.Println("非法访问: Username/Token缺失或无效!")
		userRsp.Status = -1
		userRsp.Message = "非法访问: Username/Token缺失或无效!"
		json.NewEncoder(w).Encode(userRsp)
		return
	}

	if token != r.Header.Get("Token") {
		log.Println("非法访问: Token无效或被重复登录顶号!")
		userRsp.Status = -2
		userRsp.Message = "非法访问: Token无效或被重复登录顶号!"
		json.NewEncoder(w).Encode(userRsp)
		return
	}

	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", r.Header.Get("Username")).Decode(&result)
	if len(result.Username) <= 0 {
		// 用户名不存在
		userRsp.Status = 1
		userRsp.Message = "用户名不存在!"
		json.NewEncoder(w).Encode(userRsp)
		return
	}


	var user User
	// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
	json.NewDecoder(r.Body).Decode(&user)

	if result.Photo != user.Photo {
		result.Photo = user.Photo
		// 更新  $set更新  $inc增量更新
		update := bson.D{
			{"$set", bson.D{
				{"photo", user.Photo},
			}},
		}
		models.NewMgo().UpdateOne("username", result.Username, update);
	}

	// 用户名不存在
	userRsp.Status = 0
	userRsp.Message = "修改图片成功!"
	userRsp.Result = models.PreDownload("photo", user.Photo)

	log.Println(userRsp.Message)
	json.NewEncoder(w).Encode(userRsp)

	return
}

func UserSetPassword(w http.ResponseWriter, r *http.Request) {

	var userRsp UserRsp
	var info  struct {
		Username	string `json:"username"`
		Oldpass		string `json:"oldpass"`
		Newpass		string `json:"newpass"`
	}

	log.Println(r.Body)
	// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
	json.NewDecoder(r.Body).Decode(&info)

	log.Println(info)
	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", info.Username).Decode(&result)
	if len(result.Username) <= 0 {
		// 用户名不存在
		userRsp.Status = 1
		userRsp.Message = "用户名不存在!"
		json.NewEncoder(w).Encode(userRsp)
		return
	}


	if result.Password == info.Oldpass {
		// 更新  $set更新  $inc增量更新
		update := bson.D{
			{"$set", bson.D{
				{"password", info.Newpass},
			}},
		}
		models.NewMgo().UpdateOne("username", result.Username, update);

		// 利用uuid库生成唯一且随机的token
		token := uuid.NewV4().String()
		log.Println("token", token)

		models.DeleteToken(result.Username)
		// 把token存到redis
		models.InsertToken(result.Username, token, 120);

		userRsp.Status = 0
		userRsp.Message = "修改密码成功!"
		userRsp.Token = token
	} else {
		userRsp.Status = 2
		userRsp.Message = "密码错误!"
	}

	log.Println(userRsp.Message)
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(userRsp)
	return
}

func UserSetting(w http.ResponseWriter, r *http.Request) {

	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", mux.Vars(r)["username"]).Decode(&result)

	if len(result.Username) <= 0 {
		fmt.Fprintf(w, "用户不存在!\n")
		return
	}
	log.Println(result)
	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("www/settings.tmpl")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}

	var tmplUserPage  struct {
		Username	string
		Email 		string
		Photo		string
	}
	tmplUserPage.Username = result.Username
	tmplUserPage.Email = result.Email
	if len(result.Photo) > 0 {
		tmplUserPage.Photo = models.PreDownload("photo", result.Photo)
	}

	// 利用给定数据渲染模板，并将结果写入w
	tmpl.Execute(w, tmplUserPage)
}


//[Go语言标准库之template](https://www.cnblogs.com/nickchen121/p/11517448.html)
//[GO Web编程示例 - 路由（使用gorilla/mux）](https://www.jianshu.com/p/698156c07ad4)
//[golang模板语法简明教程](https://www.cnblogs.com/Pynix/p/4154630.html)
func UserPage(w http.ResponseWriter, r *http.Request) {
	log.Println(mux.Vars(r)["username"])
	log.Println("method:", r.Method) //获取请求的方法

	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("www/user.tmpl")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}

	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", mux.Vars(r)["username"]).Decode(&result)

	var tmplUser  struct {
		Username	string
		Email 		string
		Photo		string
	}
	tmplUser.Username = result.Username
	tmplUser.Email = result.Email
	tmplUser.Photo = models.PreDownload("photo", result.Photo)

	// 利用给定数据渲染模板，并将结果写入w
	tmpl.Execute(w, tmplUser)
}

