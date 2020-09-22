package controller

import (
	"context"
	"github.com/githubchry/goweb/internal/logics"
)

func AddHandler(ctx context.Context, req *logics.AddReq) logics.AddRsp {
	// 校验token
	if logics.TokenCheck(req.Username, req.Token) != 0 {
		return logics.AddRsp{
			Code: -1,
			Message : "非法访问: Token失效!",
		}
	}

	// 校验参数
	if len(req.Operand) <= 1 {
		return logics.AddRsp{
			Code: -2,
			Message : "参数异常: 至少输入两个数!",
		}
	}

	// 调用真正的api
	return logics.Add(req)
}
