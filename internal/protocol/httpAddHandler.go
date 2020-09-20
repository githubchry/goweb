package protocol

import (
	"encoding/json"
	"github.com/githubchry/goweb/internal/controller"
	"github.com/githubchry/goweb/internal/logics"
	"log"
	"net/http"
	"strconv"
)


func HTTPAddHandler(w http.ResponseWriter, r *http.Request) {

	var req logics.AddReq
	if r.Method == "POST" {
		// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
		json.NewDecoder(r.Body).Decode(&req)
	} else if r.Method == "GET" {
		values := r.URL.Query()
		req.OperandA, _ = strconv.Atoi(values.Get("OperandA"))
		req.OperandB, _ = strconv.Atoi(values.Get("OperandB"))
	} else {
		return
	}

	log.Println("Add req: ", r.Method, req.OperandA, req.OperandB)
	rsp := controller.AddHandler(r.Context(), &req)
	log.Println("Add rsp: ", rsp)

	json.NewEncoder(w).Encode(rsp)
}

