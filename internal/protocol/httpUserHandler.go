package protocol

import (
	"github.com/githubchry/goweb/internal/controller"
	"github.com/githubchry/goweb/internal/logics"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
)

func HTTPUserLoginHandler(w http.ResponseWriter, r *http.Request) {

	req := &logics.UserLoginReq{}

	//把protobuf二进制数据转成logics.UserLoginReq结构体
	data, _ := ioutil.ReadAll(r.Body)
	if err := proto.Unmarshal(data, req); err != nil {
		log.Println("Failed to parse protobuf:", err)
		return
	}

	rsp := controller.UserLoginHandler(r.Context(), req)

	//把logics.UserLoginRsp结构体转成protobuf二进制数据
	data, _ = proto.Marshal(&rsp)
	w.Write(data)
}

func HTTPUserLogoutHandler(w http.ResponseWriter, r *http.Request) {

	req := &logics.UserLogoutReq{}

	//把protobuf二进制数据转成logics.UserLoginReq结构体
	data, _ := ioutil.ReadAll(r.Body)
	if err := proto.Unmarshal(data, req); err != nil {
		log.Println("Failed to parse protobuf:", err)
		return
	}

	controller.UserLogoutHandler(r.Context(), req)
}

func HTTPUserRegisterHandler(w http.ResponseWriter, r *http.Request) {

	req := &logics.UserRegisterReq{}

	//把protobuf二进制数据转成logics.UserLoginReq结构体
	data, _ := ioutil.ReadAll(r.Body)
	if err := proto.Unmarshal(data, req); err != nil {
		log.Println("Failed to parse protobuf:", err)
		return
	}

	rsp := controller.UserRegisterHandler(r.Context(), req)
	data, _ = proto.Marshal(&rsp)
	w.Write(data)
}

func HTTPUserSetPhotoHandler(w http.ResponseWriter, r *http.Request) {

	req := &logics.UserSetPhotoReq{}

	//把protobuf二进制数据转成logics.UserLoginReq结构体
	data, _ := ioutil.ReadAll(r.Body)
	if err := proto.Unmarshal(data, req); err != nil {
		log.Println("Failed to parse protobuf:", err)
		return
	}

	rsp := controller.UserSetPhotoHandler(r.Context(), req)
	data, _ = proto.Marshal(&rsp)
	w.Write(data)
}

func HTTPUserSetPasswordHandler(w http.ResponseWriter, r *http.Request) {

	req := &logics.UserSetPasswordReq{}

	//把protobuf二进制数据转成logics.UserLoginReq结构体
	data, _ := ioutil.ReadAll(r.Body)
	if err := proto.Unmarshal(data, req); err != nil {
		log.Println("Failed to parse protobuf:", err)
		return
	}

	rsp := controller.UserSetPasswordHandler(r.Context(), req)
	data, _ = proto.Marshal(&rsp)
	w.Write(data)
}


