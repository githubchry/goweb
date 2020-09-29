package logics

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/githubchry/goweb/internal/dao/models"
	"log"
	"net/http"
	"time"
)

/*
[](https://www.cnblogs.com/qingyunzong/p/9004509.html)
Topic
每条发布到Kafka集群的消息都有一个类别，这个类别被称为Topic。
（物理上不同Topic的消息分开存储，逻辑上一个Topic的消息虽然保存于一个或多个broker上但用户只需指定消息的Topic即可生产或消费数据而不必关心数据存于何处）
类似于数据库的表名

Partition
topic中的数据分割为一个或多个partition。每个topic至少有一个partition。
每个partition中的数据使用多个segment文件存储。partition中的数据是有序的，不同partition间的数据丢失了数据的顺序。
如果topic有多个partition，消费数据时就不能保证数据的顺序。在需要严格保证消息的消费顺序的场景下，需要将partition数目设为1。

Broker
Kafka 集群包含一个或多个服务器，服务器节点称为broker。
broker存储topic的数据(partition)。如果某topic有N个partition，集群有N个broker，那么每个broker存储该topic的一个partition。
如果某topic有N个partition，集群有(N+M)个broker，那么其中有N个broker存储该topic的一个partition，剩下的M个broker不存储该topic的partition数据。
如果某topic有N个partition，集群中broker数目少于N个，那么一个broker存储该topic的一个或多个partition。在实际生产环境中，尽量避免这种情况的发生，这种情况容易导致Kafka集群数据不均衡。



*/

type EventUploadServiceImpl struct{}

func (p *EventUploadServiceImpl) EventUpload(ctx context.Context, args *EventReq, ) (*EventRsp, error) {
	// 0.grpc收到报警事件,
	// 1.根据时间结构体imgurl字段, 使用http get取图片数据
	// 2.图片数据转为base64
	// 3.post base64到算法模块 得到特征数据
	// 4.特征数据发送到milus, 得到id
	// 5.id+特征+报警数据保存到mongo


	log.Println(args);
	rsp := &EventRsp{Message: "sucess"}
	return rsp, nil
}

var EventProducer sarama.AsyncProducer
var EventConsumer sarama.Consumer
var ResultProducer sarama.AsyncProducer
var ResultConsumer sarama.Consumer

func EventQueueInit() error {
	var err error

	EventProducer, err = models.CreateProducer()
	if err != nil {
		log.Println("NewSyncProducer", err)
		return err
	}

	EventConsumer, err = models.CreateConsumer()
	if err != nil {
		log.Println("NewConsumer", err)
		return err
	}

	ResultProducer, err = models.CreateProducer()
	if err != nil {
		log.Println("NewSyncProducer", err)
		return err
	}

	ResultConsumer, err = models.CreateConsumer()
	if err != nil {
		log.Println("NewConsumer", err)
		return err
	}

	return err
}

/*
1. 通过web api post图片到服务器
2. 服务器协程A接受 生成事件 推送到kafka 可以结构体丢topicA  图片丢topicB
3. 服务器协程B从kafka topicAB取事件, 积攒16个或超时后, 调用C++库处理(先不搞C++ 模拟出来即可)
4. 服务器从C++库得到处理结果,可能是异步回调方式,也可能是同步方式, 把结果推送到kafka topicC
5. 服务器协程C从kafka topicC获取到结果, 通过webrtc主动推到前端
*/
func EventPublish(ctx context.Context, img []byte) (*Status, error) {
	rsp := &Status{}

	// 构造一个消息
	msgstruct := &sarama.ProducerMessage{
		Topic : "event_struct",
		Value : sarama.StringEncoder("this is a test log"),
	}

	//使用通道发送
	EventProducer.Input() <- msgstruct

	timeStart := time.Now()
	//循环判断哪个通道发送过来数据.
	select {
	case sucess := <-EventProducer.Successes():
		log.Println("offset: ", sucess.Offset, "timestamp: ", sucess.Timestamp.String(), "partitions: ", sucess.Partition)
		rsp.Code = 0
		rsp.Message = "sucess"
	case fail := <-EventProducer.Errors():
		log.Println("err: ", fail.Err)
		rsp.Code = -2
		rsp.Message = "failed"
	}

	timeElapsed := time.Since(timeStart)
	log.Println("发送消息到kafka后得到结果耗时:", timeElapsed)

	return rsp, nil
}

// 通过websocket返回结果
func EventResult(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	partition := 0
	partitionConsumer, err := EventConsumer.ConsumePartition("event_struct", int32(partition), sarama.OffsetNewest)
	if err != nil {
		log.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
		return
	}
	defer partitionConsumer.Close()

	for {
		msg := <-partitionConsumer.Messages()
		log.Printf("Consumed message: [%s], offset: [%d]\n", msg.Value, msg.Offset)

		err = conn.WriteMessage(1, msg.Value)
		if err != nil {
			log.Println("write:", err)
		}
	}
}
