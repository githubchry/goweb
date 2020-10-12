package main

import (
	"github.com/githubchry/goweb/configs"
	"github.com/githubchry/goweb/internal/dao/drivers"
	"github.com/githubchry/goweb/internal/logics"
	"github.com/githubchry/goweb/internal/middleware"
	"github.com/githubchry/goweb/internal/protocol"
	"github.com/githubchry/goweb/internal/view"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
	"strconv"
)

var httpport int
var grpcport int

func initdbcfg() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间
	log.SetFlags(log.Llongfile | log.LstdFlags | log.Lmicroseconds)

	appcfg, ok := configs.LoadConfig("../configs/config.json")
	if !ok {
		return
	}
	httpport = appcfg.HTTPCfg.Port
	grpcport = appcfg.GRPCCfg.Port

	// 初始化连接到MongoDB
	err := drivers.MongoDBInit(appcfg.MongoCfg)
	if err != nil {
		log.Fatal(err)
	}

	// 初始化连接到RedisDB
	err = drivers.RedisDBInit(appcfg.RedisCfg)
	if err != nil {
		log.Fatal(err)
	}

	// 初始化连接到MinioDB
	err = drivers.MinioDBInit(appcfg.MinioCfg)
	if err != nil {
		log.Fatal(err)
	}

	// 初始化 KafkaMQ
	err = drivers.KafkaMQInit(appcfg.KafkaCfg)
	if err != nil {
		log.Fatal(err)
	}

}

func printAddr() {
	// 获取并打印一下本地ip
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				log.Printf("%s:%d\n", ipnet.IP.String(), httpport)
			}
		}
	}
}

const (
	certFileName = "../cert/cert.cer"
	keyFileName  = "../cert/cert.key"
)

func main() {

	initdbcfg()

	err := logics.EventQueueInit()
	if err != nil {
		log.Fatal(err)
	}

	// grpc
	// TLS证书解析验证
	cert, err := credentials.NewServerTLSFromFile(certFileName, keyFileName)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(cert))

	// 注册grpc方法
	logics.RegisterAddServer(grpcServer, new(logics.AddServiceImpl))
	logics.RegisterEventUploadServer(grpcServer, new(logics.EventUploadServiceImpl))

	lis, _:=net.Listen("tcp", ":"+strconv.Itoa(grpcport))
	go grpcServer.Serve(lis)

	// http2
	route := mux.NewRouter()
	route.Use(middleware.ElapsedTime)

	route.HandleFunc("/api/addpost", protocol.HTTPAddHandler) //POST
	route.HandleFunc("/api/addget", protocol.HTTPAddHandler)  //GET

	route.HandleFunc("/api/login", protocol.HTTPUserLoginHandler)                 // POST
	route.HandleFunc("/api/logout", protocol.HTTPUserLogoutHandler)               // POST
	route.HandleFunc("/api/register", protocol.HTTPUserRegisterHandler)           // POST
	route.HandleFunc("/api/userSetPhoto", protocol.HTTPUserSetPhotoHandler)       // POST
	route.HandleFunc("/api/userSetPassword", protocol.HTTPUserSetPasswordHandler) // POST
	route.HandleFunc("/api/presignedUrl", protocol.HTTPPresignedUrlHandler) // POST

	route.HandleFunc("/api/echo", logics.Echo) //WEBSOCKET

	route.HandleFunc("/user/{username}", view.HTTPUserPageHandler)            // GET
	route.HandleFunc("/settings/{username}", view.HTTPUserSettingPageHandler) // GET

	route.HandleFunc("/api/eventResult", logics.EventResult) //WEBSOCKET
	route.HandleFunc("/image/post", protocol.HTTPImagePostHandler) // GET

	route.PathPrefix("/proto").Handler(http.StripPrefix("/proto", http.FileServer(http.Dir("../proto"))))
	// 使用web目录下的文件来响应对/路径的http请求，一般用作静态文件服务，例如html、javascript、css等
	route.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../web/static"))))

	// 打印本机IP地址
	printAddr()

	// 启动http服务
	//err = http.ListenAndServeTLS(":"+strconv.Itoa(httpport), certFileName, keyFileName, route)
	err = http.ListenAndServe(":"+strconv.Itoa(httpport), route)
	if err != nil {
		log.Fatal(err)
	}

}
