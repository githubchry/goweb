package logics

func Add(req *AddReq) AddRsp {
	var rsp AddRsp
	for _, value := range req.Operand {
		rsp.Result += int64(value)
	}
	return rsp
}