package logics

import (
	"context"
	"github.com/githubchry/goweb/internal/logics/protos"
)

func Add(req *protos.AddReq) protos.AddRsp {
	var rsp protos.AddRsp
	for _, value := range req.Operand {
		rsp.Result += int64(value)
	}
	return rsp
}

type AddServiceImpl struct{}
func (p *AddServiceImpl) Add(
	ctx context.Context, args *protos.AddReq,
) (*protos.AddRsp, error) {
	reply := Add(args)
	return &reply, nil
}