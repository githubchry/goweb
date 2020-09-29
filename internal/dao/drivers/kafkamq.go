package drivers

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/githubchry/goweb/configs"
	"log"
	"strconv"
	"time"
)

var KafkaMqClient sarama.Client

func KafkaMQInit(cfg configs.KafkaCfg) error {

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


	//============ Consumer config ============
	//接收失败通知
	config.Consumer.Return.Errors = true


	var err error
	KafkaMqClient, err = sarama.NewClient([]string{cfg.Addr+":"+strconv.Itoa(cfg.Port)}, config)
	if err != nil {
		log.Fatal("create KafkaMq client failed:", err)
	}

	// 创建Topics
	err = createTopics()
	if err != nil {
		log.Fatal("createTopics Failed:", err)
	}

	//获取主题的名称集合
	topics, err := KafkaMqClient.Topics()
	if err != nil {
		log.Fatal("get topics err:", err)
	}

	for _, e := range topics {
		log.Println(e)
	}

	//获取broker集合
	brokers := KafkaMqClient.Brokers()
	//输出每个机器的地址
	for _, broker := range brokers {
		log.Println(broker.Addr())
	}

	//=================================================================================================================
	return err
}

func KafkaMQExit() {

	deleteTopics()

	KafkaMqClient.Close()
}

func createTopics() error {

	// 1.连接到Broker, 然后通过这个连接去创建Topic (很扯淡? Topic跟Broker之间不是包含关系, 但只需要理解, Broker只是提交创建Topic的请求到zookeeper, 真正创建的Topic的是zookeeper而不是其下的某一个Broker)
	// [Kafka如何创建topic](https://www.cnblogs.com/warehouse/p/9534230.html)
	// [Kafka解析之topic创建(1)](https://blog.csdn.net/u013256816/article/details/79303825)
	broker, err := KafkaMqClient.Controller()
	if err != nil {
		log.Fatal(err)
	}

	// check if the connection was OK
	connected, err := broker.Connected()
	if err != nil {
		log.Print(err.Error())
		return err
	}
	log.Print("KafkaMqClient.Controller() connect result:", connected);
	//defer KafkaMqBroker.Close()

	// 2.创建3个Topic, event_struct放结构体, 引用event_image放图片, event_status放处理结果
	topicEventStruct := "event_struct"
	topicEventStructDetail := &sarama.TopicDetail{}
	topicEventStructDetail.NumPartitions = int32(1)		//分区数
	topicEventStructDetail.ReplicationFactor = int16(1)	//副本数
	//config参数集可以用来设置topic级别的配置以覆盖默认配置。如果创建的topic再现有的集群中存在，那么会报出异常：TopicExistsException，如果创建的时候带了if-not-exists参数，那么发现topic冲突的时候可以不做任何处理
	topicEventStructDetail.ConfigEntries = make(map[string]*string)

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
	topicDetails[topicEventStruct] = topicEventStructDetail
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
	topicEventStruct := "event_struct"
	topicEventImage  := "event_image"
	topicEventStatus := "event_result"

	// 创建删除请求
	request := sarama.DeleteTopicsRequest{
		Timeout:	time.Second * 15,
		Topics: []string{topicEventStruct, topicEventImage, topicEventStatus, "event1"},
	}

	broker, err := KafkaMqClient.Controller()
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
