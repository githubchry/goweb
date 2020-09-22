package protocol

import (
	"github.com/githubchry/goweb/internal/controller"
	"github.com/githubchry/goweb/internal/logics"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)


func HTTPAddHandler(w http.ResponseWriter, r *http.Request) {

	req := &logics.AddReq{}
	if r.Method == "POST" {
		//把protobuf二进制数据转成logics.UserLoginReq结构体
		data, _ := ioutil.ReadAll(r.Body)
		if err := proto.Unmarshal(data, req); err != nil {
			log.Println("Failed to parse protobuf:", err)
			return
		}
	} else if r.Method == "GET" {

		req.Username = r.Header.Get("Username")
		req.Token = r.Header.Get("Token")

		values := r.URL.Query()
		a, _ := strconv.Atoi(values.Get("OperandA"))
		b, _ := strconv.Atoi(values.Get("OperandB"))

		req.Operand = []int32{
			int32(a),
			int32(b),
		}
	}

	rsp := controller.AddHandler(r.Context(), req)
	data, _ := proto.Marshal(&rsp)
	w.Write(data)
}

