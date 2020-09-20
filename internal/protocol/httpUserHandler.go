package protocol

import (
	"encoding/json"
	"github.com/githubchry/goweb/internal/controller"
	"github.com/githubchry/goweb/internal/logics"
	"log"
	"net/http"
)

func HTTPUserLoginHandler(w http.ResponseWriter, r *http.Request) {

	var req logics.User

	if r.Method == "POST" {
		// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
		json.NewDecoder(r.Body).Decode(&req)
	} else {
		return
	}

	log.Println("Add req: ", r.Method, req)
	rsp := controller.UserLoginHandler(r.Context(), &req)
	log.Println("Add rsp: ", rsp)

	json.NewEncoder(w).Encode(rsp)
}

func HTTPUserLogoutHandler(w http.ResponseWriter, r *http.Request) {

	controller.UserLogoutHandler(r.Context())
}

func HTTPUserRegisterHandler(w http.ResponseWriter, r *http.Request) {

	var req logics.User

	if r.Method == "POST" {
		// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
		json.NewDecoder(r.Body).Decode(&req)
	} else {
		return
	}

	log.Println("Add req: ", r.Method, req)
	rsp := controller.UserRegisterHandler(r.Context(), &req)
	log.Println("Add rsp: ", rsp)

	json.NewEncoder(w).Encode(rsp)
}

func HTTPUserSetPhotoHandler(w http.ResponseWriter, r *http.Request) {

	var req logics.User

	if r.Method == "POST" {
		// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
		json.NewDecoder(r.Body).Decode(&req)
	} else {
		return
	}

	log.Println("Add req: ", r.Method, req)
	rsp := controller.UserSetPhotoHandler(r.Context(), &req)
	log.Println("Add rsp: ", rsp)

	json.NewEncoder(w).Encode(rsp)
}

func HTTPUserSetPasswordHandler(w http.ResponseWriter, r *http.Request) {

	var req logics.UserSetPasswordRsp

	if r.Method == "POST" {
		// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
		json.NewDecoder(r.Body).Decode(&req)
	} else {
		return
	}

	log.Println("Add req: ", r.Method, req)
	rsp := controller.UserSetPasswordHandler(r.Context(), &req)
	log.Println("Add rsp: ", rsp)

	json.NewEncoder(w).Encode(rsp)
}


