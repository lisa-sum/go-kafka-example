package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	reader *kafka.Reader
	wg     sync.WaitGroup
)

func main() {
	ctx := context.Background()
	go listenSignal()
	// WriterMessage(ctx)
	// time.Sleep(5 * time.Second)
	ReaderMessage(ctx)
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

// ReaderMessage 消费消息
func ReaderMessage(ctx context.Context) {
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"192.168.0.158:9092"}, // kafka地址
		Topic:          "test-topic",                   // 要消费的topic名称
		CommitInterval: 1 * time.Second,                // 提交时间间隔
		GroupID:        "Like",                         // 消费者组名称
		StartOffset:    kafka.FirstOffset,              // 从第一条消息开始消费(仅第一次连接)
	})

	// 一直监听读取kafka消息
	for {
		if message, err := reader.ReadMessage(ctx); err != nil {
			fmt.Println("读kafka消息失败,", err)
		} else {
			fmt.Printf("%s,%d,%s,%s\n", message.Topic, message.Offset, message.Key, message.Value)
		}
	}
}

// 监听退出信号
// 当进程退出时, ReaderMessage() 会阻塞在 reader.ReadMessage(ctx) 上, 无法退出, 根本执行不到reader.Close()
// 所以需要全局定义reader, 监听退出信号, 退出时关闭reader
// syscall.SIGINT 2
// syscall.SIGTERM 15
func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	fmt.Printf("收到退出信号: %s", sig)
	if reader != nil {
		reader.Close()
	}
	os.Exit(0)
}
