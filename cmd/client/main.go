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
	rootCertFileName = "../cert/ca.pem"
	rootKeyFileName  = "../cert/ca.key"
	clientCertFileName = "../cert/client/client.pem"
	clientKeyFileName  = "../cert/client/client.key"
)

//[TLS 证书认证](https://segmentfault.com/a/1190000016601783)
func main() {
	// log打印设置: Lshortfile文件名+行号  LstdFlags日期加时间
	log.SetFlags(log.Llongfile | log.LstdFlags)

	// TLS证书解析验证
	cert, err := tls.LoadX509KeyPair(clientCertFileName, clientKeyFileName)
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile(rootCertFileName)
	if err != nil {
		log.Fatal(err)
	}

	if true != certPool.AppendCertsFromPEM(ca) {
		log.Fatal(err)
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},	//客户端证书
		ServerName: 	"chry-server",
		RootCAs:      	certPool,
	})

	log.Println("-------------------------")
	var conn *grpc.ClientConn
	if true {
		conn, err = grpc.Dial("127.0.0.1:8848", grpc.WithTransportCredentials(creds), grpc.WithBlock())
	} else {
		conn, err = grpc.Dial("127.0.0.1:8848", grpc.WithInsecure(), grpc.WithBlock())
	}
	log.Println("-------------------------")
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
