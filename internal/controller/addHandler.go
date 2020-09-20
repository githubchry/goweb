package controller

import (
	"context"
	"github.com/githubchry/goweb/internal/logics"
)

func AddHandler(ctx context.Context, req *logics.AddReq) logics.AddRsp {
	// 校验token
	username, _ := ctx.Value("username").(string)
	token, _ := ctx.Value("token").(string)
	if logics.TokenCheck(username, token) != 0 {
		var rsp logics.AddRsp
		rsp.Status.Code = -2
		rsp.Status.Message = "非法访问: Token失效!"
		return rsp
	}

	// 校验参数 < int32


	// 调用真正的api
	return logics.Add(req)
}
