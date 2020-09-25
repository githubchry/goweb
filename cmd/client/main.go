package main

import (
	"context"
	"github.com/githubchry/goweb/internal/logics"
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

	conn, err := grpc.Dial("127.0.0.1:7070", grpc.WithTransportCredentials(cert), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := logics.NewAddClient(conn)

	req := &logics.AddReq{}
	req.Operand = []int32{123, 456,}

	reply, err := client.Add(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply.Result)
}
