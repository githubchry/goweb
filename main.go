package main

import (
	"github.com/deepch/vdk/format/rtsp"
	"github.com/githubchry/goweb/drivers"
	"github.com/githubchry/goweb/webapi"
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

func init()  {
	// 初始化连接到MongoDB
	err := drivers.MongoDBInit()
	if err != nil {
		log.Fatal(err)
		return
	}

	// 初始化连接到RedisDB
	err = drivers.RedisDBInit()
	if err != nil {
		log.Fatal(err)
		return
	}

	// 初始化连接到MinioDB
	err = drivers.MinioDBInit()
	if err != nil {
		log.Fatal(err)
		return
	}
}


func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间  LstdFlags  [Go语言标准库之log](https://www.cnblogs.com/nickchen121/p/11517450.html)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// rtsp 转 webrtc
	rtsp.DebugRtsp = true // 打印rtsp流程

	route := mux.NewRouter()

	// 处理发往/api/add的http请求
	route.HandleFunc("/api/addpost", 		webapi.Addpost)   //POST
	route.HandleFunc("/api/addget", 		webapi.Addget)     //GET
	route.HandleFunc("/api/getcodec", 		webapi.Getcodec) //POST
	route.HandleFunc("/api/swapsdp", 		webapi.Swapsdp)   //POST
	route.HandleFunc("/api/login", 		webapi.UserLogin)   		// POST
	route.HandleFunc("/api/logout", 		webapi.UserLogout) 	// POST
	route.HandleFunc("/api/register", 		webapi.UserRegister)	// POST
	route.HandleFunc("/api/userSetPhoto", 		webapi.UserSetPhoto)	// POST
	route.HandleFunc("/api/userSetPassword", 		webapi.UserSetPassword)	// POST


	route.HandleFunc("/api/presignedUrl", 	webapi.PresignedUrl)	// POST
	route.HandleFunc("/api/echo", 			webapi.Echo)             //WEBSOCKET

	route.HandleFunc("/user/{username}", 		webapi.UserPage)   		// POST
	route.HandleFunc("/settings/{username}", 		webapi.UserSetting)   		// POST

	// 使用web目录下的文件来响应对/路径的http请求，一般用作静态文件服务，例如html、javascript、css等
	route.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("www/"))))

	// 打印本机IP地址
	printAddr()

	// 启动http服务
	err := http.ListenAndServe(":"+port, route)
	if err != nil {
		log.Fatal(err)
	}

}
