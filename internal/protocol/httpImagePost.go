package protocol

import (
	"encoding/json"
	"github.com/githubchry/goweb/internal/logics"
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

	rsp, _:= logics.ImagePostHandler(r.Context(), buf)


	json_str, err := json.Marshal(rsp)
	log.Printf("%s\n", json_str)


	w.Write(json_str)
}


