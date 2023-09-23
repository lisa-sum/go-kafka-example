package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

func main() {
	ctx := context.Background()
	WriterMessage(ctx)
}

// WriterMessage 生产消息
func WriterMessage(ctx context.Context) {
	topic := "test-topic" // topic名称
	writer := kafka.Writer{
		Addr:                   kafka.TCP("192.168.0.158:9092"), // kafka地址
		Balancer:               &kafka.Hash{},                   // 负载均衡算法,用于分配topic到哪个消费者
		WriteTimeout:           3 * time.Second,                 // 写时间超时设置
		RequiredAcks:           kafka.RequireNone,               // 使用哪种方式对per返回消息成功的处理
		AllowAutoTopicCreation: true,                            // 是否允许自动创建Topic
	}
	defer writer.Close()

	// 允许重试的次数
	for i := 0; i < 3; i++ {
		if err := writer.WriteMessages(
			ctx,
			kafka.Message{Topic: topic, Key: []byte("1"), Value: []byte("Hello1 1")},
			kafka.Message{Topic: topic, Key: []byte("2"), Value: []byte("Hello2 1")},
			kafka.Message{Topic: topic, Key: []byte("3"), Value: []byte("Hello3 1")},
			kafka.Message{Topic: topic, Key: []byte("4"), Value: []byte("Hello4 1")},
			kafka.Message{Topic: topic, Key: []byte("5"), Value: []byte("Hello5 1")},
		); err != nil {
			// 第一次没有topic, 通常都是失败的
			if err == kafka.LeaderNotAvailable {
				time.Sleep(500 * time.Millisecond)
			} else {
				fmt.Printf("批量写入kafka失败!, 如果你是第一次写入topic, 通常都是失败的%s", err)
			}
		} else {
			break
		}
	}
}
