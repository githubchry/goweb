package controller

import (
	"context"
	"github.com/githubchry/goweb/internal/logics"
)

func PresignedUrlHandler(ctx context.Context, req *logics.FileReq) logics.FileRsp {
	// 校验token
	username, _ := ctx.Value("username").(string)
	token, _ := ctx.Value("token").(string)
	if logics.TokenCheck(username, token) != 0 {
		var rsp logics.FileRsp
		rsp.Status.Code = -2
		rsp.Status.Message = "非法访问: Token失效!"
		return rsp
	}

	// 校验参数


	// 调用真正的api
	return logics.PresignedUrl(*req)
}

