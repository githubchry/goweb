package drivers

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDbConn *mongo.Client
var MongoDbName string


// 连接
func connect() (*mongo.Client, error) {

	// 设置客户端参数
	clientOptions := options.Client().ApplyURI("mongodb://chry:chry@localhost:27017/?authSource=test")

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	//defer MongoDbConn.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	return client, err
}

// 初始化
func MongoDBInit() error{
	var err error
	MongoDbConn, err = connect()
	if err != nil {
		log.Fatal(err)
	}
	// 不要用 defer MongoDbConn.Disconnect(context.TODO())
	MongoDbName = "test"
	return err
}

// 关闭
func MongoDBExit() {
	err := MongoDbConn.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("MongoDB is closed.")
}
