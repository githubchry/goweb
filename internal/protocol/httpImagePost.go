package protocol

import (
	"encoding/json"
	"github.com/githubchry/goweb/internal/logics"
	"io/ioutil"
	"log"
	"net/http"
)

func HTTPImagePostHandler(w http.ResponseWriter, r *http.Request) {
	// 根据字段名获取表单文件
	formFile, _, err := r.FormFile("image")
	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		return
	}
	defer formFile.Close()

	data, _ := ioutil.ReadAll(formFile)

	rsp, err := logics.ImagePostHandler(r.Context(), data)
	json_str, err := json.Marshal(rsp)
	log.Printf("%s\n", json_str)

	w.Write(json_str)
}


