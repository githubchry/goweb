package controller

import (
	"context"
	"github.com/githubchry/goweb/internal/logics"
	"github.com/githubchry/goweb/internal/logics/protos"
)

func PresignedUrlHandler(ctx context.Context, req *protos.FileReq) protos.FileRsp {
	// 校验token
	if logics.TokenCheck(req.Username, req.Token) != 0 {
		return protos.FileRsp{
			Code: -1,
			Message : "非法访问: Token失效!",
		}
	}

	// 校验参数
	if req.Cmd != 0 && req.Cmd != 1 {
		return protos.FileRsp{
			Code: -2,
			Message : "参数异常: Cmd错误!",
		}
	}

	if req.Type != "photo" && req.Type != "music" {
		return protos.FileRsp{
			Code: -2,
			Message : "参数异常: Type仅支持'photo'或'music'!",
		}
	}

	if req.Filename == "" {
		return protos.FileRsp{
			Code: -2,
			Message : "参数异常: Name不能为空!",
		}
	}

	// 调用真正的api
	return logics.PresignedUrl(req)
}

