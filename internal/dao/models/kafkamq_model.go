package models

import (
	"github.com/Shopify/sarama"
	"github.com/githubchry/goweb/internal/dao/drivers"
	"log"
)



func loopConsumer(consumer sarama.Consumer, topic string, partition int) {
	partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
	if err != nil {
		log.Println(err)
		return
	}
	defer partitionConsumer.Close()

	for {
		msg := <-partitionConsumer.Messages()
		log.Printf("Consumed message: [%s], offset: [%d]\n", msg.Value, msg.Offset)
	}
}

func CreateProducer() (sarama.AsyncProducer, error) {

	//使用配置,新建一个异步生产者
	producer, err := sarama.NewAsyncProducerFromClient(drivers.KafkaMqClient)
	if err != nil {
		log.Println("NewSyncProducer", err)
	}

	return producer, err
}

func CreateConsumer() (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumerFromClient(drivers.KafkaMqClient)
	if err != nil {
		log.Println("NewSyncProducer", err)
	}
	return consumer, err
}