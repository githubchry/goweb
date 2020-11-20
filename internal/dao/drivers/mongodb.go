package drivers

import (
	"context"
	"github.com/githubchry/goweb/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
)

var MongoDbConn *mongo.Client
var MongoDbName string

// 初始化
func MongoDBInit(cfg configs.MongoCfg) error {
	var err error

	MongoDbName = cfg.DBName
	// 设置客户端参数	"mongodb://chry:chry@localhost:27017/?authSource=test"
	url := "mongodb://" + cfg.Username + ":" + cfg.Password + "@" + cfg.Addr + ":" + strconv.Itoa(cfg.Port) + "/?authSource=" + cfg.DBName
	clientOptions := options.Client().ApplyURI(url)

	// 连接到MongoDB
	MongoDbConn, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = MongoDbConn.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	// 不要用 defer MongoDbConn.Disconnect(context.TODO())
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
