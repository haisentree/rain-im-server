package global

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

const (
	// 客户端连接不在网关服务器上，就投放在消息队列中，由消息队列分发给网关服务器
	GatewayRelayMessageTheme = "gateway.message.relay.%s"
)

const (
	GatewayStreamNameMessage       = "GATEWAY-PERSIST-MESSAGE"
	GatewayMessageSeqIncreaseTheme = "gateway.message.seq.increase.%s"
	GatewayMessageSeqIncreaseQueue = "gateway-message-seq-increase" //  这个队列中单个消费者顺序消费
	GatewayMessageSaveTheme        = "gateway.message.save.%s"
	GatewayMessageSaveQueue        = "gateway-message-save"
)

func initJetStream(js nats.JetStreamContext) error {
	subjects := []string{
		fmt.Sprintf(GatewayMessageSeqIncreaseTheme, "*"), // seq + 1
		fmt.Sprintf(GatewayMessageSaveTheme, "*"),        // save message
	}

	_, err := js.StreamInfo(GatewayStreamNameMessage)
	if err == nil {
		return err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      GatewayStreamNameMessage,
		Subjects:  subjects,
		Storage:   nats.FileStorage,
		Retention: nats.WorkQueuePolicy,
	})
	if err != nil {
		log.Fatalf("创建流失败: %v", err)
	}

	return err
}

// // 订阅主题
// sub, err := nc.Subscribe("foo", func(msg *nats.Msg) {
// 	log.Printf("收到消息: %s", string(msg.Data))
// })
// if err != nil {
// 	log.Fatal(err)
// }
// defer sub.Unsubscribe()

// // 发布消息
// err = nc.Publish("foo", []byte("Hello NATS!"))
// if err != nil {
// 	log.Fatal(err)
// }

// // 持久化
// js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
// if err != nil {
// 	log.Fatal(err)
// }

// // 定义流配置
// stream, err := js.AddStream(&nats.StreamConfig{
// 	Name:      "ORDERS-2",
// 	Subjects:  []string{"orders-2.>"}, // 匹配 orders. 开头的所有主题
// 	Storage:   nats.FileStorage,       // 持久化到磁盘
// 	Retention: nats.WorkQueuePolicy,
// })
// if err != nil {
// 	fmt.Println("js2")
// 	log.Fatal(err)
// }

// fmt.Println("stream:", stream)

// // 发布消息到流
// ack, err := js.Publish("orders-2.new", []byte("order-1"))
// if err != nil {
// 	log.Fatal(err)
// }
// log.Printf("消息存储于 stream: %s, seq: %d", ack.Stream, ack.Sequence)

// // 创建消费者（自动创建）
// sub2, err := js.QueueSubscribe("orders-2.>", "worker-group", func(msg *nats.Msg) {
// 	log.Printf("收到订单: %s", string(msg.Data))
// 	msg.Ack()
// }, nats.ManualAck())
// if err != nil {
// 	log.Fatal(err)
// }
// defer sub2.Unsubscribe()
