package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type AddReq struct {
	OperandA int
	OperandB int
}

type AddReply struct {
	Result int
}

func add(w http.ResponseWriter, r *http.Request) {
	var addReq AddReq
	// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
	json.NewDecoder(r.Body).Decode(&addReq)
	// 打印请求数据
	log.Println("req: ", addReq.OperandA, addReq.OperandB)
	var addReply AddReply
	// 进行加法计算，并保存结果到结构体内
	addReply.Result = addReq.OperandA + addReq.OperandB
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(addReply)
}

//websocket由http升级而来，首先会发送附带Upgrade请求头的Http请求，所以我们需要在处理Http请求时拦截请求并判断其是否为websocket升级请求
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 放行跨域请求
	},
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
func main() {
	// 获取并打印一下本地ip
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Printf("%s:8080\n", ipnet.IP.String())
			}
		}
	}

	// 处理发往/api/add的http请求
	http.HandleFunc("/api/add", add)
	http.HandleFunc("/echo", echo)

	// 使用web目录下的文件来响应对/路径的http请求，一般用作静态文件服务，例如html、javascript、css等
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web/"))))

	// 启动http服务
	log.Fatal(http.ListenAndServe(":8080", nil))

}
