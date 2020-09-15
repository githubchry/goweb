package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"github.com/githubchry/goweb/drivers"
	"github.com/githubchry/goweb/models"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/codec/h264parser"
	"github.com/deepch/vdk/format/rtsp"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

// AddReq represents the parameter of an addition operation.
type AddReq struct {
	OperandA int
	OperandB int
}

// AddRsp represents the result of an addition operation.
type AddRsp struct {
	Result int		`json:"result"`	// 结果 OperandA + OperandB
	Status int		`json:"status"`	// 状态 0表示sucess
	Message string	`json:"message"`// 消息
}

func addpost(w http.ResponseWriter, r *http.Request) {

	//log.Println(r.Header.Get("Username"), r.Header.Get("Token"))

	var addRsp AddRsp
	// 从redis查询token
	token, err := models.FindToken(r.Header.Get("Username"))
	if err != nil {
		log.Println("非法访问: Username/Token缺失或无效!")
		addRsp.Status = -1
		addRsp.Message = "非法访问: Username/Token缺失或无效!"
		json.NewEncoder(w).Encode(addRsp)
		return
	}

	if token != r.Header.Get("Token") {
		log.Println("非法访问: Token无效或被重复登录顶号!")
		addRsp.Status = -2
		addRsp.Message = "非法访问: Token无效或被重复登录顶号!"
		json.NewEncoder(w).Encode(addRsp)
		return
	}

	// 延长token生存周期?

	var addReq AddReq
	// 将请求的body作为JSON字符串解码，并存入AddReq结构体内
	json.NewDecoder(r.Body).Decode(&addReq)
	// 打印请求数据
	log.Println("post req: ", addReq.OperandA, addReq.OperandB)
	// 进行加法计算，并保存结果到结构体内
	addRsp.Result = addReq.OperandA + addReq.OperandB
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(addRsp)
}

func addget(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Header.Get("Username"), r.Header.Get("Token"))

	var addRsp AddRsp
	// 从redis查询token
	token, err := models.FindToken(r.Header.Get("Username"))
	if err != nil {
		log.Println("非法访问: Username/Token缺失或无效!")
		addRsp.Status = -1
		addRsp.Message = "非法访问: Username/Token缺失或无效!"
		json.NewEncoder(w).Encode(addRsp)
		return
	}

	if token != r.Header.Get("Token") {
		log.Println("非法访问: Token无效或被重复登录顶号!")
		addRsp.Status = -2
		addRsp.Message = "非法访问: Token无效或被重复登录顶号!"
		json.NewEncoder(w).Encode(addRsp)
		return
	}

	values := r.URL.Query()

	var addReq AddReq
	addReq.OperandA, _ = strconv.Atoi(values.Get("OperandA"))
	addReq.OperandB, _ = strconv.Atoi(values.Get("OperandB"))

	// 打印请求数据
	log.Println("get req: ", addReq.OperandA, addReq.OperandB)

	// 进行加法计算，并保存结果到结构体内
	addRsp.Result = addReq.OperandA + addReq.OperandB
	// 将结果结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(addRsp)
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


func user(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //获取请求的方法
	// 解析url传递的参数
	r.ParseForm()
	for k, v := range r.Form {
		log.Println("key:", k)
		// join() 方法用于把数组中的所有元素放入一个字符串。
		// 元素是通过指定的分隔符进行分隔的
		log.Println("val:", strings.Join(v, ""))
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

var codec []av.CodecData
var session *rtsp.Client

func getcodec(w http.ResponseWriter, r *http.Request) {
	rtspurl, _ := ioutil.ReadAll(r.Body)
	// 打印请求数据
	log.Println("post req: ", string(rtspurl))

	// 创建session
	var err error
	session, err = rtsp.Dial(string(rtspurl))
	if err != nil {
		log.Println(err)
		return
	}
	// 设置session rtp超时时间为10s
	session.RtpKeepAliveTimeout = 10 * time.Second

	// 获取rtsp流 OPTIONS->DESCRIBE->SETUP->PLAY
	codec, err = session.Streams()
	if err != nil {
		log.Println(err)
		return
	}
	// 将codec结构体进行JSON编码，并写入响应
	json.NewEncoder(w).Encode(codec)

}

func swapsdp(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	websdpbase64 := string(body)
	// 打印请求数据
	log.Println("post websdp: ", websdpbase64)

	websdp, err := base64.StdEncoding.DecodeString(websdpbase64)
	if err != nil {
		log.Println("DecodeString error", err)
		return
	}

	// Create Media MediaEngine
	mediaEngine := webrtc.MediaEngine{}
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  string(websdp),
	}
	err = mediaEngine.PopulateFromSDP(offer)
	if err != nil {
		log.Println("PopulateFromSDP error", err)
		return
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))
	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Println("NewPeerConnection error", err)
		return
	}

	// ADD Video Track
	var payloadType uint8
	payloadType = 102
	videoTrack, err := peerConnection.NewTrack(payloadType, rand.Uint32(), "video", "chry")
	if err != nil {
		log.Fatalln("NewTrack", err)
	}
	_, err = peerConnection.AddTransceiverFromTrack(videoTrack,
		webrtc.RtpTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionSendonly,
		},
	)
	if err != nil {
		log.Println("AddTransceiverFromTrack error", err)
		return
	}
	_, err = peerConnection.AddTrack(videoTrack)
	if err != nil {
		log.Println("AddTrack error", err)
		return
	}

	// peerConnection设置对端sdp
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		log.Println("SetRemoteDescription error", err, offer.SDP)
		return
	}
	// 创建本地应答sdp
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Println("CreateAnswer error", err)
		return
	}
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		log.Println("SetLocalDescription error", err)
		return
	}

	log.Println("")
	log.Println("")
	// 返回本地sdp
	answersdp := base64.StdEncoding.EncodeToString([]byte(answer.SDP))
	log.Println("answersdp: ", answersdp)
	fmt.Fprintf(w, "%s\n", answersdp)

	// 至此两者已经交换好sdp

	// ADD KeepAlive Timer
	// web端创建了DataChannel并每秒发送'ping', 这里进行接收处理
	timer1 := time.NewTimer(time.Second * 2)
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
			timer1.Reset(2 * time.Second)
		})
	})

	control := make(chan bool, 10) // 通道控制
	// 设置连接状态改变回调函数 检测到连接成功时进行发送媒体数据处理
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("Connection State has changed %s \n", connectionState.String())
		if connectionState != webrtc.ICEConnectionStateConnected {
			// 连接不成功直接退出
			log.Println("Client Exit Exit")
			err := peerConnection.Close()
			if err != nil {
				log.Println("peerConnection Exit error", err)
			}
			control <- true
			return
		}

		// 至此说明连接成功 起一条协程往对端发送媒体数据
		go func() {

			// 获取sps pps
			if codec == nil {
				log.Println("Codec error")
				return
			}
			sps := codec[0].(h264parser.CodecData).SPS()
			pps := codec[0].(h264parser.CodecData).PPS()

			var start bool // 发送标志
			var Vts time.Duration
			var Vpre time.Duration

			for {
				// 从流里面读数据  注意拿到的nalu前4个字节是长度, 后面的data是没有00 00 00 01的
				pkt, err := session.ReadPacket()
				if err != nil {
					log.Println(err)
					break
				}

				log.Printf("Idx%d, %2x \n", pkt.Idx, pkt.Data[0:10])

				if pkt.IsKeyFrame {
					// 收到I帧才开始发送
					start = true
				}

				// pkt.Idx为0表示视频帧 pkt.Idx为1表示音频帧 感觉说法有点扯淡应该不是固定的 但老子是新手 先抄下来这么着吧
				// 只发送视频帧
				if false == start || pkt.Idx != 0 {
					continue
				}

				if pkt.IsKeyFrame {
					pkt.Data = append([]byte("\000\000\001"+string(sps)+"\000\000\001"+string(pps)+"\000\000\001"), pkt.Data[4:]...)

				} else {
					pkt.Data = append([]byte("\000\000\001"), pkt.Data[4:]...)
				}

				if Vpre != 0 {
					Vts = pkt.Time - Vpre
				}
				samples := uint32(90000 / 1000 * Vts.Milliseconds())
				log.Println("samples", samples, pkt.Time)
				err = videoTrack.WriteSample(media.Sample{Data: pkt.Data, Samples: samples})
				if err != nil {
					return
				}
				Vpre = pkt.Time

			}
			log.Println("3333333333333333")

			err = session.Close()
			if err != nil {
				log.Println("session Exit error", err)
			}
			log.Println("reconnect wait 5s")

		}()

	})
}

func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间  LstdFlags  [Go语言标准库之log](https://www.cnblogs.com/nickchen121/p/11517450.html)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

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

	// 查询总数
	name, size := models.NewMgo().Count()
	log.Printf(" documents name: %+v documents size %d \n", name, size)

	// rtsp 转 webrtc
	rtsp.DebugRtsp = true // 打印rtsp流程

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
	http.HandleFunc("/api/addpost", addpost)   //POST
	http.HandleFunc("/api/addget", addget)     //GET
	http.HandleFunc("/api/getcodec", getcodec) //POST
	http.HandleFunc("/api/swapsdp", swapsdp)   //POST
	http.HandleFunc("/api/login", login)   		// POST
	http.HandleFunc("/api/logout", logout) 	// POST
	http.HandleFunc("/api/register", register)	// POST
	http.HandleFunc("/api/echo", echo)             //WEBSOCKET

	// 使用web目录下的文件来响应对/路径的http请求，一般用作静态文件服务，例如html、javascript、css等
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web/"))))

	// 启动http服务
	log.Fatal(http.ListenAndServe(":8080", nil))

	// 断开连接
	drivers.MongoDBExit()
	drivers.RedisDBExit()
}
