package protocol

import (
	"encoding/json"
	"github.com/githubchry/goweb/internal/logics"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

func HTTPImagePostHandler(w http.ResponseWriter, r *http.Request) {
	var rsp *logics.AlgorithmOutput
	// 根据字段名获取表单文件
	formFile, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		return
	}
	defer formFile.Close()

	data, _ := ioutil.ReadAll(formFile)

	if len(data) > logics.MAXMSGBYTES - 20 {
		log.Printf("图片过大!!")

		rsp = &logics.AlgorithmOutput{
			Code: -1,
			Message: "图片过大!",
		}

		log.Printf("图片过大!")
		json_str, _ := json.Marshal(rsp)
		log.Printf("%s\n", json_str)

		w.Write(json_str)
		return
	}

	ext := path.Ext(header.Filename)
	filetype := logics.ImageFile_FILE_TYPE_UNKNOW

	if ext == ".jpg" ||  ext == ".jpeg"{
		filetype = logics.ImageFile_FILE_TYPE_JPG
	} else if ext == ".png" {
		filetype = logics.ImageFile_FILE_TYPE_PNG
	}

	log.Printf(ext)

	var req logics.AlgorithmInput
	req.Images = []*logics.ImageFile{
		{
			Data: data,
			Type: filetype,
		},
	}

	rsp, err = logics.ImagePostHandler(r.Context(), &req)
	json_str, err := json.Marshal(rsp)
	log.Printf("%s\n", json_str)

	w.Write(json_str)
}


