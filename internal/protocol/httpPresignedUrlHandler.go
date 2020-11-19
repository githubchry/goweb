package protocol

import (
	"github.com/githubchry/goweb/internal/controller"
	"github.com/githubchry/goweb/internal/logics/protos"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
)

func HTTPPresignedUrlHandler(w http.ResponseWriter, r *http.Request) {

	req := &protos.FileReq{}

	//把protobuf二进制数据转成logics.UserLoginReq结构体
	data, _ := ioutil.ReadAll(r.Body)
	if err := proto.Unmarshal(data, req); err != nil {
		log.Println("Failed to parse protobuf:", err)
		return
	}

	rsp := controller.PresignedUrlHandler(r.Context(), req)
	data, _ = proto.Marshal(&rsp)
	w.Write(data)
}

