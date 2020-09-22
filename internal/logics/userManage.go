package logics

import (
	"github.com/githubchry/goweb/internal/dao/models"
	"github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)


// [golang jwt-go的使用](https://www.cnblogs.com/jianga/p/12487267.html)
// [使用JWT进行接口认证](https://studygolang.com/articles/27242?fr=sidebar)
func UserLogin(req *UserLoginReq) UserLoginRsp {

	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", req.Username).Decode(&result)

	var rsp UserLoginRsp

	if result.Password == req.Password {
		rsp.Token = TokenGenerate(result.Username)
		rsp.Message = "登录成功!"
		rsp.Code = 0
	} else if len(result.Username) <= 0 {
		rsp.Code = 1
		rsp.Message = "用户名不存在!"
	} else {
		rsp.Code = 2
		rsp.Message = "密码错误!"
	}

	log.Println(rsp.Message)
	return rsp
}

func UserLogout(req *UserLogoutReq) {

	log.Println(req)
	// 从redis查询token
	ret, err := models.FindToken(req.Username)
	if err != nil {
		return
	}

	if req.Token != ret {
		return
	}

	TokenDelete(req.Username)
	log.Println("注销用户:", req.Username)
	return
}

func UserRegister(req *UserRegisterReq) Status {

	// 打印请求数据
	log.Println("post req: ", req.Username, req.Password, req.Email)

	// 操作mongodb
	// 返回结果
	var rsp Status

	// 删除
	//deleteResult := models.NewMgo().DeleteMany("username", req.Username)
	//fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult)

	// 查询用户名是否已存在
	var result User
	models.NewMgo().FindOne("username", req.Username).Decode(&result)
	if len(result.Username) > 0 {
		// 用户名已注册
		rsp.Code = 1
		rsp.Message = "用户名已注册!"
	} else {
		// 单个插入
		InsertOneResult := models.NewMgo().InsertOne(req)
		log.Println("Inserted a single document: ", InsertOneResult)

		rsp.Code = 0
		rsp.Message = "注册成功!"
	}

	log.Println(rsp.Message)
	return rsp
}

func UserSetPhoto(req *UserSetPhotoReq) Status {

	var userRsp Status
	// 查询用户名是否已存在
	var result User
	_ = models.NewMgo().FindOne("username", req.Username).Decode(&result)
	if len(result.Username) <= 0 {
		// 用户名不存在
		userRsp.Code = 1
		userRsp.Message = "用户名不存在!"
		return userRsp
	}

	if result.Photo != req.Photo {
		result.Photo = req.Photo
		// 更新  $set更新  $inc增量更新
		update := bson.D{
			{"$set", bson.D{
				{"photo", req.Photo},
			}},
		}
		models.NewMgo().UpdateOne("username", result.Username, update);
	}

	userRsp.Code = 0
	userRsp.Message = "修改图片成功!"
	log.Println(userRsp.Message)

	return userRsp
}

func UserSetPassword(req *UserSetPasswordReq) Status {

	log.Println(req)
	// 查询用户名是否已存在
	var result User
	var userRsp Status
	models.NewMgo().FindOne("username", req.Username).Decode(&result)
	if len(result.Username) <= 0 {
		// 用户名不存在
		userRsp.Code = 1
		userRsp.Message = "用户名不存在!"
		return userRsp
	}

	if result.Password == req.Oldpass {
		// 更新  $set更新  $inc增量更新
		update := bson.D{
			{"$set", bson.D{
				{"password", req.Newpass},
			}},
		}
		models.NewMgo().UpdateOne("username", result.Username, update);

		// 利用uuid库生成唯一且随机的token
		token := uuid.NewV4().String()
		log.Println("token", token)

		models.DeleteToken(result.Username)

		userRsp.Code = 0
		userRsp.Message = "修改密码成功!"
	} else {
		userRsp.Code = 2
		userRsp.Message = "密码错误!"
	}

	log.Println(userRsp.Message)
	return userRsp
}
