package main

import (
	"context"
	"github.com/githubchry/goweb/internal/logics/protos"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)


const (
	certFileName = "../cert/cert.cer"
)

//[TLS 证书认证](https://segmentfault.com/a/1190000016601783)
func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间
	log.SetFlags(log.Llongfile | log.LstdFlags)
	// TLS证书解析验证

	cert, err := credentials.NewClientTLSFromFile(certFileName, "chry-server")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial("127.0.0.1:7070", grpc.WithTransportCredentials(cert))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// 模拟报警事件上传
	client := protos.NewEventUploadClient(conn)
	event := &protos.EventReq{
		Time:   ptypes.TimestampNow(),
		Type:   protos.EventReq_EVENT_TYPE_SUSPECT,
		Addr:   "192.168.1.99",
		Token:  "",
		Imgurl: "",
	}

	reply, err := client.EventUpload(context.Background(), event)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(reply.Message)
}
