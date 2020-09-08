package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type AddReq struct {
	OperandA int
	OperandB int
}

type AddReply struct {
	Result int
}

func addpost(w http.ResponseWriter, r *http.Request) {
	var addReq AddReq
	// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
	json.NewDecoder(r.Body).Decode(&addReq)
	// 打印请求数据
	log.Println("post req: ", addReq.OperandA, addReq.OperandB)
	var addReply AddReply
	// 进行加法计算，并保存结果到结构体内
	addReply.Result = addReq.OperandA + addReq.OperandB
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(addReply)
}

func addget(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	var addReq AddReq
	addReq.OperandA, _ = strconv.Atoi(values.Get("OperandA"))
	addReq.OperandB, _ = strconv.Atoi(values.Get("OperandB"))

	// 打印请求数据
	log.Println("get req: ", addReq.OperandA, addReq.OperandB)

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

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("web/login.html")
		t.Execute(w, nil)
	} else {
		//请求的是登录数据，那么执行登录的逻辑判断
		//解析传过来的参数，默认不会解析，必须显示调用后服务器才会输出参数信息
		r.ParseForm() //解析form
		//这里的request.Form["username"]可以用request.FormValue("username")代替，那么就不需要显示调用  request.ParseForm
		fmt.Printf("username: %v\n", r.Form["username"])
		fmt.Printf("password: %v\n", r.Form["password"])
	}
}

func user(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	// 解析url传递的参数
	r.ParseForm()
	for k, v := range r.Form {
		fmt.Println("key:", k)
		// join() 方法用于把数组中的所有元素放入一个字符串。
		// 元素是通过指定的分隔符进行分隔的
		fmt.Println("val:", strings.Join(v, ""))
	}
	// 输出到客户端
	name := r.Form["username"]
	pass := r.Form["password"]
	for _, v := range name {
		fmt.Fprintf(w, "用户名:%v\n", v)
	}
	for _, n := range pass {
		fmt.Fprintf(w, "密码:%v\n", n)
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
	http.HandleFunc("/api/addpost", addpost) //POST
	http.HandleFunc("/api/addget", addget)   //GET
	http.HandleFunc("/login", login)         //GET + POST
	http.HandleFunc("/user", user)           //GET + POST
	http.HandleFunc("/echo", echo)

	// 使用web目录下的文件来响应对/路径的http请求，一般用作静态文件服务，例如html、javascript、css等
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web/"))))

	// 启动http服务
	log.Fatal(http.ListenAndServe(":8080", nil))

}
