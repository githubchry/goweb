package controller

import (
	"context"
	"github.com/githubchry/goweb/internal/logics"
)

func UserLoginHandler(ctx context.Context, req *logics.User) logics.UserRsp {
	var rsp logics.UserRsp
	// 无需校验token

	// 校验参数
	if req.Username == "" {
		return rsp
	}

	if req.Password == "" {
		return rsp
	}

	// 调用真正的api
	return logics.UserLogin(*req)
}


func UserLogoutHandler(ctx context.Context) {
	// 校验token
	username, _ := ctx.Value("username").(string)
	token, _ := ctx.Value("token").(string)
	if logics.TokenCheck(username, token) != 0 {
		return
	}

	// 无需校验参数

	// 调用真正的api
	logics.UserLogout(username, token)
}

func UserRegisterHandler(ctx context.Context, req *logics.User) logics.UserRsp {
	var rsp logics.UserRsp
	// 无需校验token

	// 校验参数
	if req.Username == "" {
		return rsp
	}

	if req.Email == "" {
		return rsp
	}

	if req.Password == "" {
		return rsp
	}

	// 调用真正的api
	return logics.UserRegister(*req)
}

func UserSetPhotoHandler(ctx context.Context, req *logics.User) logics.UserRsp {
	var rsp logics.UserRsp
	// 无需校验token

	// 校验参数
	if req.Username == "" {
		return rsp
	}

	if req.Photo == "" {
		return rsp
	}

	// 调用真正的api
	return logics.UserSetPhoto(*req)
}

func UserSetPasswordHandler(ctx context.Context, req *logics.UserSetPasswordRsp) logics.UserRsp {
	var rsp logics.UserRsp
	// 无需校验token

	// 校验参数
	if req.Username == "" {
		return rsp
	}

	if req.Newpass == "" {
		return rsp
	}

	if req.Oldpass == "" {
		return rsp
	}
	// 调用真正的api
	return logics.UserSetPassword(*req)
}
