package webapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/codec/h264parser"
	"github.com/deepch/vdk/format/rtsp"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var codec []av.CodecData
var session *rtsp.Client

func Getcodec(w http.ResponseWriter, r *http.Request) {
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

func Swapsdp(w http.ResponseWriter, r *http.Request) {
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

