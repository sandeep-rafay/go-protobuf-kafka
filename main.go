package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	v1 "go-protobuf/gen/proto/go/proto/v1"
)

func main() {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:19092", "compression.type": "gzip"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "product"

	p1 := &v1.Product{
		Id:    100001,
		Name:  "Shirt",
		Brand: "Acme",
	}

	p2 := &v1.Product{
		Id:    100002,
		Name:  "T-Shirt",
		Brand: "Acme",
	}

	p3 := &v1.Product{
		Id:    100003,
		Name:  "Shorts",
		Brand: "Acme",
	}

	for _, product := range []*v1.Product{p1, p2, p3} {

		id, _ := uuid.New().MarshalBinary()
		mProduct, _ := proto.Marshal(product)

		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            id,
			Value:          mProduct,
		}, nil)
		if err != nil {
			return
		}
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
