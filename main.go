package main

import (
	"github.com/deepch/vdk/format/rtsp"
	"github.com/githubchry/goweb/drivers"
	"github.com/githubchry/goweb/models"
	"github.com/githubchry/goweb/webapi"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间  LstdFlags  [Go语言标准库之log](https://www.cnblogs.com/nickchen121/p/11517450.html)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// rtsp 转 webrtc
	rtsp.DebugRtsp = true // 打印rtsp流程

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

	// 查询总数
	name, size := models.NewMgo().Count()
	log.Printf(" documents name: %+v documents size %d \n", name, size)

	// 获取并打印一下本地ip
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				log.Printf("%s:8080\n", ipnet.IP.String())
			}
		}
	}

	// 处理发往/api/add的http请求
	http.HandleFunc("/api/addpost", webapi.Addpost)   //POST
	http.HandleFunc("/api/addget", webapi.Addget)     //GET
	http.HandleFunc("/api/getcodec", webapi.Getcodec) //POST
	http.HandleFunc("/api/swapsdp", webapi.Swapsdp)   //POST
	http.HandleFunc("/api/login", webapi.UserLogin)   		// POST
	http.HandleFunc("/api/logout", webapi.UserLogout) 	// POST
	http.HandleFunc("/api/register", webapi.UserRegister)	// POST
	http.HandleFunc("/api/presignedUrl", webapi.PresignedUrl)	// POST

	http.HandleFunc("/api/echo", webapi.Echo)             //WEBSOCKET

	// 使用web目录下的文件来响应对/路径的http请求，一般用作静态文件服务，例如html、javascript、css等
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./www/"))))

	// 启动http服务
	log.Fatal(http.ListenAndServe(":8080", nil))

	// 断开连接
	drivers.MongoDBExit()
	drivers.RedisDBExit()
}
