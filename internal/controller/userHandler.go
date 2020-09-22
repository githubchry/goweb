package controller

import (
	"context"
	"github.com/githubchry/goweb/internal/logics"
	"log"
	"regexp"
)

func UserLoginHandler(ctx context.Context, req *logics.UserLoginReq) logics.UserLoginRsp {

	var rsp logics.UserLoginRsp
	// 校验参数
	//用户名以字母下划线开头，由数字和字母组成 2-16位
	match, _ := regexp.MatchString("^[a-zA-z_]\\w{1,15}$", req.Username)
	if match == false {
		log.Println(match)
		rsp.Code = 1;
		rsp.Message = "用户名格式错误";
		return rsp
	}

	if len(req.Password) != 32 {
		rsp.Code = 2;
		rsp.Message = "密码md5异常错误";
		return rsp
	}

	// 调用真正的api
	return logics.UserLogin(req)
}


func UserLogoutHandler(ctx context.Context, req *logics.UserLogoutReq) {
	// 调用真正的api
	logics.UserLogout(req)
}

func UserRegisterHandler(ctx context.Context, req *logics.UserRegisterReq) logics.Status {
	var rsp logics.Status
	// 无需校验token

	// 校验参数
	//用户名以字母下划线开头，由数字和字母组成 2-16位
	match, _ := regexp.MatchString("^[a-zA-z_]\\w{1,15}$", req.Username)
	if match == false {
		rsp.Code = 1;
		rsp.Message = "用户名格式错误";
		return rsp
	}

	//电子邮箱 前缀由字母、数字、下划线、短线“-”、点号“.”组成，后缀域名由字母、数字、短线“-”、域名后缀组成   ;
	match, _ = regexp.MatchString("^(\\w-*\\.*)+@(\\w-?)+(\\.\\w{2,})+$", req.Email)
	if match == false {
		rsp.Code = 1;
		rsp.Message = "电子邮箱格式错误";
		return rsp
	}

	if len(req.Password) != 32 {
		rsp.Code = 1;
		rsp.Message = "密码md5格式错误";
		return rsp
	}

	// 调用真正的api
	return logics.UserRegister(req)
}

func UserSetPhotoHandler(ctx context.Context, req *logics.UserSetPhotoReq) logics.Status {
	var rsp logics.Status
	// 无需校验token

	// 校验参数
	//用户名以字母下划线开头，由数字和字母组成 2-16位
	match, _ := regexp.MatchString("^[a-zA-z_]\\w{1,15}$", req.Username)
	if match == false {
		rsp.Code = 1;
		rsp.Message = "用户名格式错误";
		return rsp
	}

	// 头像图片以jpg|png结尾
	match, _ = regexp.MatchString(".*(.jpg|.png)$", req.Photo)
	if match == false {
		rsp.Code = 2;
		rsp.Message = "图片格式错误";
		return rsp
	}

	// 调用真正的api
	return logics.UserSetPhoto(req)
}

func UserSetPasswordHandler(ctx context.Context, req *logics.UserSetPasswordReq) logics.Status {
	var rsp logics.Status
	// 无需校验token

	// 校验参数
	//用户名以字母下划线开头，由数字和字母组成 2-16位
	match, _ := regexp.MatchString("^[a-zA-z_]\\w{1,15}$", req.Username)
	if match == false {
		rsp.Code = 1;
		rsp.Message = "用户名格式错误";
		return rsp
	}

	if len(req.Newpass) != 32 || len(req.Oldpass) != 32 {
		rsp.Code = 2;
		rsp.Message = "密码md5异常错误";
		return rsp
	}

	// 调用真正的api
	return logics.UserSetPassword(req)
}
