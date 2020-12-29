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

const TopicImage = "event_image"
const TopicStatus = "event_result"

func KafkaMQInit(cfg configs.KafkaCfg) error {

	//详细的config参数:https://blog.csdn.net/chinawangfei/article/details/93097203
	config := sarama.NewConfig()
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V2_6_0_0
	config.Net.DialTimeout = 300000000	//3秒 不要用3 * time.Second 对应不上! => 0.3*time.Second

	KafkaMqAddr = []string{cfg.Addr+":"+strconv.Itoa(cfg.Port)}
	var err error
	log.Println("KafkaMQ Client Conn .....")
	kafkaMqClient, err = sarama.NewClient(KafkaMqAddr, config)
	if err != nil {
		log.Println("create KafkaMq", KafkaMqAddr, "client failed:", err)
		return err
	}

	// 创建Topics
	err = createTopics()
	if err != nil {
		log.Println("createTopics Failed:", err)
		return err
	}

	//获取主题的名称集合
	topics, err := kafkaMqClient.Topics()
	if err != nil {
		log.Println("get topics err:", err)
		return err
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
	cleanup_policy := "delete"
	retention_bytes := "40000000"	//400M
	segment_bytes 	:= "20000000"			//200M

	configEntries := map[string]*string{
		"cleanup.policy" : &cleanup_policy,
		"retention.bytes": &retention_bytes,
		"segment.bytes": &segment_bytes,
	}

	topicEventImageDetail := &sarama.TopicDetail{}
	topicEventImageDetail.NumPartitions = int32(1)
	topicEventImageDetail.ReplicationFactor = int16(1)
	topicEventImageDetail.ConfigEntries = configEntries

	topicEventStatusDetail := &sarama.TopicDetail{}
	topicEventStatusDetail.NumPartitions = int32(1)
	topicEventStatusDetail.ReplicationFactor = int16(1)
	topicEventStatusDetail.ConfigEntries = configEntries

	//message_max_bytes := "100000"
	//topicEventImageDetail.ConfigEntries["message.max.bytes"] = &message_max_bytes

	topicDetails := make(map[string]*sarama.TopicDetail)
	topicDetails[TopicImage] = topicEventImageDetail
	topicDetails[TopicStatus] = topicEventStatusDetail

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
		} else {
			log.Printf("topic [%s] 创建成功!\n", key)
		}
	}

	return nil
}

func deleteTopics() {
	// 创建删除请求
	request := sarama.DeleteTopicsRequest{
		Timeout:	time.Second * 15,
		Topics: []string{TopicImage, TopicStatus},
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


