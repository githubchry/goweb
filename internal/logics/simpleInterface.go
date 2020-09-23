package logics

import "context"

func Add(req *AddReq) AddRsp {
	var rsp AddRsp
	for _, value := range req.Operand {
		rsp.Result += int64(value)
	}
	return rsp
}

type AddServiceImpl struct{}
func (p *AddServiceImpl) Add(
	ctx context.Context, args *AddReq,
) (*AddRsp, error) {
	reply := Add(args)
	return &reply, nil
}