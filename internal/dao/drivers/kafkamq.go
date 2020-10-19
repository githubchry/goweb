package drivers

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/githubchry/goweb/configs"
	"log"
	"strconv"
	"time"
)

var kafkaMqClient sarama.Client
var KafkaMqAddr []string

func KafkaMQInit(cfg configs.KafkaCfg) error {

	//详细的config参数:https://blog.csdn.net/chinawangfei/article/details/93097203
	config := sarama.NewConfig()
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V2_6_0_0

	KafkaMqAddr = []string{cfg.Addr+":"+strconv.Itoa(cfg.Port)}
	var err error
	kafkaMqClient, err = sarama.NewClient(KafkaMqAddr, config)
	if err != nil {
		log.Fatal("create KafkaMq client failed:", err)
	}

	// 创建Topics
	err = createTopics()
	if err != nil {
		log.Fatal("createTopics Failed:", err)
	}

	//获取主题的名称集合
	topics, err := kafkaMqClient.Topics()
	if err != nil {
		log.Fatal("get topics err:", err)
	}

	for _, e := range topics {
		log.Println(e)
	}

	//获取broker集合
	brokers := kafkaMqClient.Brokers()
	//输出每个机器的地址
	for _, broker := range brokers {
		log.Println(broker.Addr())
	}

	//=================================================================================================================
	return err
}

func KafkaMQExit() {

	deleteTopics()

	kafkaMqClient.Close()
}

func createTopics() error {

	// 1.连接到Broker, 然后通过这个连接去创建Topic (很扯淡? Topic跟Broker之间不是包含关系, 但只需要理解, Broker只是提交创建Topic的请求到zookeeper, 真正创建的Topic的是zookeeper而不是其下的某一个Broker)
	// [Kafka如何创建topic](https://www.cnblogs.com/warehouse/p/9534230.html)
	// [Kafka解析之topic创建(1)](https://blog.csdn.net/u013256816/article/details/79303825)
	broker, err := kafkaMqClient.Controller()
	if err != nil {
		log.Fatal(err)
	}

	// check if the connection was OK
	connected, err := broker.Connected()
	if err != nil {
		log.Print(err.Error())
		return err
	}
	log.Print("kafkaMqClient.Controller() connect result:", connected);
	//defer KafkaMqBroker.Close()

	// 2.创建2个Topic, event_image放图片, event_status放处理结果
	topicEventImage := "event_image"
	topicEventImageDetail := &sarama.TopicDetail{}
	topicEventImageDetail.NumPartitions = int32(1)
	topicEventImageDetail.ReplicationFactor = int16(1)
	topicEventImageDetail.ConfigEntries = make(map[string]*string)

	topicEventStatus := "event_result"
	topicEventStatusDetail := &sarama.TopicDetail{}
	topicEventStatusDetail.NumPartitions = int32(1)
	topicEventStatusDetail.ReplicationFactor = int16(1)
	topicEventStatusDetail.ConfigEntries = make(map[string]*string)

	topicDetails := make(map[string]*sarama.TopicDetail)
	topicDetails[topicEventImage] = topicEventImageDetail
	topicDetails[topicEventStatus] = topicEventStatusDetail

	// 创建请求
	request := sarama.CreateTopicsRequest{
		Timeout:      time.Second * 15,
		TopicDetails: topicDetails,
	}

	// 发送请求
	response, err := broker.CreateTopics(&request)
	if err != nil {
		log.Printf("%#v", &err)
		return err
	}

	for key, val := range response.TopicErrors {
		if val.Err != 0 {
			if val.Err == 36 {
				log.Printf("topic [%s] 已存在, 无需重复创建!\n", key)
			} else {
				log.Printf("create topic [%s] error: %s\n", key, val.Error())
				return errors.New(val.Error())
			}
		}
	}

	return nil
}

func deleteTopics() {
	topicEventImage  := "event_image"
	topicEventStatus := "event_result"

	// 创建删除请求
	request := sarama.DeleteTopicsRequest{
		Timeout:	time.Second * 15,
		Topics: []string{topicEventImage, topicEventStatus},
	}

	broker, err := kafkaMqClient.Controller()
	if err != nil {
		log.Println(err)
		return
	}

	response, err := broker.DeleteTopics(&request)
	if err != nil {
		log.Printf("%#v", &err)
	}

	for key, val := range response.TopicErrorCodes {
		if val != 0 {
			log.Printf("delete topic [%s] error: %s\n", key, val.Error())
		}
	}
}
