package protocol

import (
	"encoding/json"
	"github.com/githubchry/goweb/internal/controller"
	"github.com/githubchry/goweb/internal/logics"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
)

func HTTPUserLoginHandler(w http.ResponseWriter, r *http.Request) {

	var req logics.UserLoginReq

	if r.Method == "POST" {
		//把protobuf二进制数据转成logics.UserLoginReq结构体
		data, _ := ioutil.ReadAll(r.Body)
		if err := proto.Unmarshal(data, &req); err != nil {
			log.Fatalln("Failed to parse r.Body:", err)
		}
	} else {
		return
	}

	rsp := controller.UserLoginHandler(r.Context(), &req)

	//把logics.UserLoginRsp结构体转成protobuf二进制数据
	data, _ := proto.Marshal(&rsp)
	n, _ := w.Write(data)
	log.Println(len(data), "/", n, ":",data)
}

func HTTPUserLogoutHandler(w http.ResponseWriter, r *http.Request) {

	controller.UserLogoutHandler(r.Context())
}

func HTTPUserRegisterHandler(w http.ResponseWriter, r *http.Request) {

	var req logics.UserRegisterReq

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

	var req logics.UserSetPhotoReq

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

	var req logics.UserSetPasswordReq

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


