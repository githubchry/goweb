package protocol

import (
	"encoding/json"
	"github.com/githubchry/goweb/internal/controller"
	"github.com/githubchry/goweb/internal/logics"
	"log"
	"net/http"
)

func HTTPPresignedUrlHandler(w http.ResponseWriter, r *http.Request) {

	var req logics.FileReq

	if r.Method == "POST" {
		// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
		json.NewDecoder(r.Body).Decode(&req)
	} else {
		return
	}

	log.Println("Add req: ", r.Method, req)
	rsp := controller.PresignedUrlHandler(r.Context(), &req)
	log.Println("Add rsp: ", rsp)

	json.NewEncoder(w).Encode(rsp)
}

