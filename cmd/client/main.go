package main

import (
	"context"
	"github.com/githubchry/goweb/internal/logics"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:1234", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := logics.NewAddClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &logics.AddReq{}
	req.Operand = []int32{
		123,
		4156,
	}

	reply, err := client.Add(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply.Result)
}
