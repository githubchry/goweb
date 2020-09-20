package main

import (
	"github.com/githubchry/goweb/internal/logics"
	"github.com/githubchry/goweb/internal/middleware"
	"github.com/githubchry/goweb/internal/protocol"
	"github.com/githubchry/goweb/internal/view"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

const port = "8080"

func printAddr(){
	// 获取并打印一下本地ip
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				log.Printf("%s:%s\n", ipnet.IP.String(), port)
			}
		}
	}
}


func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间  LstdFlags
	log.SetFlags(log.Llongfile | log.LstdFlags)

	route := mux.NewRouter()
	route.Use(middleware.ElapsedTime)
	route.Use(middleware.ReadToken)

	route.HandleFunc("/api/addpost", 			protocol.HTTPAddHandler)         //POST
	route.HandleFunc("/api/addget", 			protocol.HTTPAddHandler)         //GET

	route.HandleFunc("/api/login", 			protocol.HTTPUserLoginHandler)      // POST
	route.HandleFunc("/api/logout", 			protocol.HTTPUserLogoutHandler)    // POST
	route.HandleFunc("/api/register", 		protocol.HTTPUserRegisterHandler)           // POST
	route.HandleFunc("/api/userSetPhoto", 	protocol.HTTPUserSetPhotoHandler)       // POST
	route.HandleFunc("/api/userSetPassword", 	protocol.HTTPUserSetPasswordHandler) // POST

	route.HandleFunc("/api/presignedUrl", 	protocol.HTTPPresignedUrlHandler) // POST

	route.HandleFunc("/api/echo", 			logics.Echo)                 //WEBSOCKET

	route.HandleFunc("/user/{username}",	 	view.HTTPUserPageHandler)        // GET
	route.HandleFunc("/settings/{username}", 	view.HTTPUserSettingPageHandler) // GET

	// 使用web目录下的文件来响应对/路径的http请求，一般用作静态文件服务，例如html、javascript、css等
	route.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../web/static"))))

	// 打印本机IP地址
	printAddr()

	// 启动http服务
	err := http.ListenAndServe(":"+port, route)
	if err != nil {
		log.Fatal(err)
	}

}
