package logics

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/githubchry/goweb/internal/dao/models"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

type EventUploadServiceImpl struct{}
func (p *EventUploadServiceImpl) EventUpload(ctx context.Context, args *EventReq, ) (*EventRsp, error) {
	log.Println(args);
	rsp := &EventRsp{Message: "sucess"}
	return rsp, nil
}

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
const MAXMSGBYTES = 1000000

var EventProducer sarama.AsyncProducer
var ResultProducer sarama.AsyncProducer

var EventImageConsumer sarama.Consumer
var EventStructConsumer sarama.Consumer
var ResultConsumer sarama.Consumer

var wslist []*websocket.Conn
var wschan chan []byte
var msgchan chan EventReq	// 暂时没啥卵用
var imgchan chan []byte


func consumerThread(topic string, process func(msg *sarama.ConsumerMessage)) {
	msgchan = make(chan EventReq, 16)

	partitionList, err := EventStructConsumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		log.Printf("fail to get list of %v partition:%v\n", topic, err)
	}

	// partition号从0开始
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		partitionConsumer, err := EventStructConsumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("failed to start consumer for %v partition %d,err:%v\n", topic, partition, err)
			return
		}
		defer partitionConsumer.AsyncClose()

		// 异步从每个分区消费信息
		func(sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				//log.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
				process(msg)
			}
		}(partitionConsumer)
	}
}

func eventStructConsumer(msg *sarama.ConsumerMessage) {
	//反序列化后 去获取image  ->chan
	var req EventReq
	err := proto.Unmarshal(msg.Value, &req)
	if err != nil {
		log.Printf("failed to Unmarshal proto:%v\n", err)
		return
	}
	msgchan <- req
}

func eventImageConsumer(msg *sarama.ConsumerMessage) {
	imgchan <- msg.Value
}

func resultConsumer(msg *sarama.ConsumerMessage) {
	// 把消息发送到websocket通道
	wschan <- msg.Value
}


func eventHandle() {

	for {
		select {
		case img := <- imgchan:
			log.Printf("图片大小:%v\n", len(img))
			strresult := "已处理图片: size="+strconv.Itoa(len(img))
			log.Printf(strresult)

			// 转化成kafka消息
			msgresult := &sarama.ProducerMessage{
				Topic : "event_result",
				Value : sarama.StringEncoder(strresult),
			}

			err := ProducerInput(ResultProducer, msgresult)
			if err != nil {
				log.Println("eventProducer struct failed:", err)
				return
			}
		default:
			log.Println("=================该批次已处理完毕!=================")
			return
		}
	}
}

func eventHandlerThread() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		imgnum := len(imgchan)
		select {
		case <-ticker.C:
			if imgnum > 0 {
				log.Printf("单位时间内未囤满事件:%v/16\n", imgnum)
				eventHandle()
				imgnum = 0
			}

		case msg := <- msgchan:
			log.Printf("收到事件:%v\n", msg.Offset)
		default:
			if imgnum >= 16 {
				log.Printf("单位时间内囤满16个事件\n")
				eventHandle()
				imgnum = 0
				ticker.Reset(10 * time.Second)
			}
		}
	}
}

// 创建Queue的4个成员: 事件生产/消费者 结果消费/生产者
func QueueMemberInit() error {
	var err error

	//详细的config参数:https://blog.csdn.net/chinawangfei/article/details/93097203
	config := sarama.NewConfig()
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V2_6_0_0

	//============ Producer config ============
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//等待服务器所有副本都保存成功后的响应 	发送完数据需要leader和follow都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse(默认)这里才有用. 成功交付的消息将在success channel返回
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.MaxMessageBytes = MAXMSGBYTES	// 保持默认

	EventProducer, err = models.CreateProducer(*config)
	if err != nil {
		log.Println("EventProducer", err)
		return err
	}

	ResultProducer, err = models.CreateProducer(*config)
	if err != nil {
		log.Println("NewSyncProducer", err)
		return err
	}

	//============ Consumer config ============
	//接收失败通知
	config.Consumer.Return.Errors = true
	ResultConsumer, err = models.CreateConsumer(*config)
	if err != nil {
		log.Println("NewConsumer", err)
		return err
	}

	EventStructConsumer, err = models.CreateConsumer(*config)
	if err != nil {
		log.Println("NewConsumer", err)
		return err
	}


	//config.Consumer.Fetch.Max = 16
	//config.Consumer.Fetch.Min = 16
	//config.Consumer.MaxWaitTime = time.Second * 10
	EventImageConsumer, err = models.CreateConsumer(*config)
	if err != nil {
		log.Println("NewConsumer", err)
		return err
	}
	return nil
}

func EventQueueInit() error {
	var err error

	err = QueueMemberInit()
	if err != nil {
		log.Println("QueueMemberInit", err)
		return err
	}
	msgchan = make(chan EventReq, 16)
	imgchan = make(chan []byte, 16)
	// 创建事件结果消费线程
	wschan = make(chan []byte)
	go wspolling()

	go consumerThread("event_struct", eventStructConsumer)
	go consumerThread("event_image", eventImageConsumer)
	go consumerThread("event_result", resultConsumer)

	go eventHandlerThread()
	return err
}

// 监听最新结果, 发送到wslist上的每一个ws
func wspolling() {
	for {
		log.Println("websocket轮询线程")
		select {
		case msg := <-wschan:
			log.Println("websocket轮询线程收到事件处理结果:", string(msg))
			// 轮询转发过程中移除已关闭的客户端 [slice移除算法出处](https://blog.csdn.net/liyunlong41/article/details/85132603)
			idx := 0	// 记录下一个有效conn应该在的位置
			for _, conn := range wslist{
				log.Println("转发消息至:", &conn)

				// 判断客户端是否已经断开
				err := conn.WriteMessage(1, msg)
				if err != nil {
					log.Println("检测到客户端[", &conn, "]已断开")
					conn.Close()
				} else {
					// 正常连接的客户端 更新到wslist[idx]
					wslist[idx] = conn
					idx++
				}
			}
			// for过程中已经更新了list, 对结果进行截取
			wslist = wslist[:idx]
		}
	}
}

func ProducerInput(producer sarama.AsyncProducer, msg *sarama.ProducerMessage) error {

	//使用通道发送
	producer.Input() <- msg
	//循环判断哪个通道发送过来数据.
	select {
	case sucess := <-producer.Successes():
		log.Println("offset: ", sucess.Offset, "timestamp: ", sucess.Timestamp.String(), "partitions: ", sucess.Partition)
		return nil
	case fail := <-producer.Errors():
		log.Println("err: ", fail.Err)
		return fail.Err
	}
}

/*
1. 通过web api post图片到服务器
2. 服务器协程A接受 生成事件 推送到kafka 可以结构体丢topicA  图片丢topicB
3. 服务器协程B从kafka topicAB取事件, 积攒16个或超时后, 调用C++库处理(先不搞C++ 模拟出来即可)
4. 服务器从C++库得到处理结果,可能是异步回调方式,也可能是同步方式, 把结果推送到kafka topicC
5. 服务器协程C从kafka topicC获取到结果, 通过webrtc主动推到前端
*/

func ImagePostHandler(ctx context.Context, img []byte) (*Status, error) {
	rsp := &Status{Message: "已上传!"}

	if len(img) > MAXMSGBYTES {
		rsp.Code = -1
		rsp.Message = "图片过大!"
		return rsp, nil
	}

	// 先发送image, 然后发送struct
	msgimage := &sarama.ProducerMessage{
		Topic : "event_image",
		Value : sarama.ByteEncoder(img),
	}

	err := ProducerInput(EventProducer, msgimage)
	if err != nil {
		log.Println("eventProducer image failed:", err)
		rsp.Code = -1
		rsp.Message = "fail"
		return rsp, err
	}

	// 随便构造一个报警消息
	event := &EventReq{
		Time: ptypes.TimestampNow(),
		Type: EventReq_EVENT_TYPE_SUSPECT,
		Addr: "192.168.1.99",
		Token: "",
		Imgurl: "",
		Offset: msgimage.Offset,
	}
	data, _ := proto.Marshal(event)

	// 转化成kafka消息
	msgstruct := &sarama.ProducerMessage{
		Topic : "event_struct",
		Value : sarama.ByteEncoder(data),
	}

	err = ProducerInput(EventProducer, msgstruct)
	if err != nil {
		log.Println("ProducerInput struct failed:", err)
		rsp.Code = -2
		rsp.Message = "fail"
		return rsp, err
	}

	return rsp, nil
}

// 通过websocket返回结果
func EventResult(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	wslist = append(wslist, conn)
	log.Println("当前websocket客户端个数:", len(wslist))
}
