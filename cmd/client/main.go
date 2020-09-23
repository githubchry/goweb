package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/githubchry/goweb/internal/logics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

const (
	certFile = "../cert/cert.pem"
	keyFile =  "../cert/key.pem"
)

//[TLS 证书认证](https://segmentfault.com/a/1190000016601783)
func main() {

	// TLS证书解析验证
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},//客户端证书
		ServerName: 	"chry-server",
		RootCAs:      	certPool,
	})

	log.Println("-------------------------")
	var conn *grpc.ClientConn
	if true {
		conn, err = grpc.Dial("127.0.0.1:8848", grpc.WithTransportCredentials(creds))
	} else {
		conn, err = grpc.Dial("127.0.0.1:8848", grpc.WithInsecure(), grpc.WithBlock())
	}
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := logics.NewAddClient(conn)

	req := &logics.AddReq{}
	req.Operand = []int32{
		123,
		456,
	}

	reply, err := client.Add(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply.Result)
}
