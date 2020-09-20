package logics

// api请求结果
type Status struct {
	Code 	int		`json:"code"`	// 状态码 0表示sucess
	Message string	`json:"message"`// 消息
}
