package models

import (
	"github.com/Shopify/sarama"
	"github.com/githubchry/goweb/internal/dao/drivers"
	"log"
)


func LoopConsumer(consumer sarama.Consumer, topic string, process func(msg *sarama.ConsumerMessage)) {
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		log.Printf("fail to get list of %v partition:%v\n", topic, err)
	}

	// partition号从0开始
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("failed to start consumer for %v partition %d,err:%v\n", topic, partition, err)
			return
		}
		defer partitionConsumer.AsyncClose()

		// 异步从每个分区消费信息
		func(sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				//log.Printf("Partition:%d Offset:%d Key:%v len(Value):%v", msg.Partition, msg.Offset, msg.Key, len(msg.Value))
				process(msg)
			}
		}(partitionConsumer)
	}
}

func CreateConsumer(config sarama.Config) (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumer(drivers.KafkaMqAddr, &config)
	if err != nil {
		log.Println("NewSyncProducer", err)
	}
	return consumer, err
}

func CreateProducer(config sarama.Config) (sarama.AsyncProducer, error) {
	//使用配置,新建一个异步生产者
	producer, err := sarama.NewAsyncProducer(drivers.KafkaMqAddr, &config)
	if err != nil {
		log.Println("NewSyncProducer", err)
	}

	return producer, err
}

func ProducerInput(producer sarama.AsyncProducer, msg *sarama.ProducerMessage) (int64, error) {
	//使用通道发送
	producer.Input() <- msg
	//循环判断哪个通道发送过来数据.
	select {
	case sucess := <-producer.Successes():
		//log.Println("offset: ", sucess.Offset, "timestamp: ", sucess.Timestamp.String(), "partitions: ", sucess.Partition)
		return sucess.Offset, nil
	case fail := <-producer.Errors():
		log.Println("err: ", fail.Err)
		return -1, fail.Err
	}
}
