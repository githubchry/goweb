package logics

import (
	"github.com/githubchry/goweb/internal/dao/models"
	"github.com/satori/go.uuid"
	"log"
)

func TokenGenerate(username string) string {
	// 利用uuid库生成唯一且随机的token
	token := uuid.NewV4().String()
	log.Println("token", token)

	// 把token存到redis
	models.InsertToken(username, token, 120);
	return token
}

func TokenCheck(username string, token string) int {
	// 从redis查询token
	ret, err := models.FindToken(username)
	if err != nil {
		return -1
	}

	if token != ret {
		return -2
	}
	return 0
}

func TokenDelete(username string) {
	models.DeleteToken(username)
}

