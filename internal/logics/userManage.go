package logics

import (
	"github.com/githubchry/goweb/internal/dao/models"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

type User struct {
	Username	string 	`json:"username"`
	Password	string	`json:"password"`
	Email 		string	`json:"email"`
	Photo 		string	`json:"photo"`
}


// UserRsp represents the result of an addition operation.
type UserRsp struct {
	Status 	Status	`json:"status"`			//应答状态
	Result 	string	`json:"result"`
	Token 	string 	`json:"token"`
}

type UserSetPasswordRsp  struct {
	Username	string `json:"username"`
	Oldpass		string `json:"oldpass"`
	Newpass		string `json:"newpass"`
}
// [golang jwt-go的使用](https://www.cnblogs.com/jianga/p/12487267.html)
// [使用JWT进行接口认证](https://studygolang.com/articles/27242?fr=sidebar)
func UserLogin(user User) UserRsp {
	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", user.Username).Decode(&result)

	var rsp UserRsp
	if result.Password == user.Password {
		rsp.Token = TokenGenerate(result.Username)
		rsp.Status.Code = 0
		rsp.Status.Message = "登录成功!"
	} else if len(result.Username) <= 0 {
		rsp.Status.Code = 1
		rsp.Status.Message = "用户名不存在!"
	} else {
		rsp.Status.Code = 2
		rsp.Status.Message = "密码错误!"
	}

	log.Println(rsp.Status.Message)
	return rsp
}

func UserLogout(username string, token string) {

	// 从redis查询token
	ret, err := models.FindToken(username)
	if err != nil {
		return
	}

	if token != ret {
		return
	}

	TokenDelete(username)
	log.Println("注销用户:", username)
	return
}

func UserRegister(user User) UserRsp {

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
		rsp.Status.Code = 1
		rsp.Status.Message = "用户名已注册!"
	} else {
		// 单个插入
		InsertOneResult := models.NewMgo().InsertOne(user)
		log.Println("Inserted a single document: ", InsertOneResult)

		rsp.Status.Code = 0
		rsp.Status.Message = "注册成功!"
	}

	log.Println(rsp.Status.Message)
	return rsp
}

func UserSetPhoto(user User) UserRsp {

	var userRsp UserRsp
	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", user.Username).Decode(&result)
	if len(result.Username) <= 0 {
		// 用户名不存在
		userRsp.Status.Code = 1
		userRsp.Status.Message = "用户名不存在!"
		return userRsp
	}

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

	userRsp.Status.Code = 0
	userRsp.Status.Message = "修改图片成功!"
	userRsp.Result = models.PreDownload("photo", user.Photo)

	log.Println(userRsp.Status.Message)

	return userRsp
}

func UserSetPassword(info UserSetPasswordRsp) UserRsp {

	log.Println(info)
	// 查询用户名是否已存在
	var result User
	var userRsp UserRsp
	models.NewMgo().FindOne("username", info.Username).Decode(&result)
	if len(result.Username) <= 0 {
		// 用户名不存在
		userRsp.Status.Code = 1
		userRsp.Status.Message = "用户名不存在!"
		return userRsp
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

		userRsp.Status.Code = 0
		userRsp.Status.Message = "修改密码成功!"
		userRsp.Token = token
	} else {
		userRsp.Status.Code = 2
		userRsp.Status.Message = "密码错误!"
	}

	log.Println(userRsp.Status.Message)
	return userRsp
}
