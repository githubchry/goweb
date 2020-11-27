package protocol

import (
	"github.com/githubchry/goweb/internal/logics"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
)

func HTTPImagePostHandler(w http.ResponseWriter, r *http.Request) {

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("method:", r.Method, "buf len:", len(buf)) //获取请求的方法

	rsp, _ := logics.ImagePostHandler(r.Context(), buf)
	data, _ := proto.Marshal(rsp)
	w.Write(data)
}
