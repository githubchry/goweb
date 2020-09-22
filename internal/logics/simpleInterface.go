package logics

// AddReq represents the parameter of an addition operation.
type AddReq struct {
	OperandA int	`json:"OperandA" validate:"required||integer=-65535, 65535"`
	OperandB int	`json:"OperandB"`
}

// AddRsp represents the result of an addition operation.
type AddRsp struct {
	Result 	int		`json:"result"`	// 结果 OperandA + OperandB
	Status 	Status	`json:"status"`	// 状态 0表示sucess
}

func Add(req *AddReq) AddRsp {
	var rsp AddRsp
	rsp.Result = req.OperandA + req.OperandB
	return rsp
}