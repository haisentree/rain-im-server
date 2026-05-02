package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"rain-im-server/global"
	gatewayv1 "rain-im-server/proto/protogo/gateway/v1"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/encoding/protojson"
)

type ClientWorker struct {
	subscriptions []*nats.Subscription
}

func NewClientWorker() *MessageWorker {
	return &MessageWorker{}
}

func (w *ClientWorker) Run() {
	w.Subscribe()

	select {}
}

// 优雅关闭
func (w *ClientWorker) Stop() {
	for _, sub := range w.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			log.Printf("取消订阅失败: %v", err)
		}
	}
	log.Println("消息工作器已停止")
}

func (w *ClientWorker) Subscribe() {
	// 如果需要同时订阅 Seq 增加队列，可以添加
	seqSubject := fmt.Sprintf("%s.*", global.GatewayMessageSeqIncreaseTheme[:len(global.GatewayMessageSeqIncreaseTheme)-3])
	seqSub, err := global.NatsJS.QueueSubscribe(
		seqSubject,
		"message-seq-worker-group",
		func(msg *nats.Msg) {
			w.handleMessage(msg)
		},
		nats.ManualAck(),
		nats.Durable("message-seq-consumer"),
	)
	if err != nil {
		log.Fatalf("订阅消息序号队列失败: %v", err)
	}

	w.subscriptions = append(w.subscriptions, seqSub)
	log.Printf("消息序号工作器已启动，订阅主题: %s", seqSubject)

}

// 处理消息持久化
func (w *ClientWorker) handleMessage(msg *nats.Msg) {
	defer msg.Ack()

	// 1. 解析原始消息
	var rawMsg gatewayv1.RawMessage
	if err := json.Unmarshal(msg.Data, &rawMsg); err != nil {
		log.Printf("解析 RawMessage 失败: %v", err)
		return
	}

	// 2. 解析 SingleMessage
	var singleMsg gatewayv1.SingleMessage
	if err := protojson.Unmarshal(rawMsg.Data, &singleMsg); err != nil {
		log.Printf("解析 SingleMessage 失败: %v", err)
		return
	}
}

// 处理序号增加
func (w *MessageWorker) handleSeqIncrease(msg *nats.Msg) {
	defer msg.Ack()

	var rawMsg gatewayv1.RawMessage
	if err := json.Unmarshal(msg.Data, &rawMsg); err != nil {
		log.Printf("解析 RawMessage 失败: %v", err)
		return
	}

	var singleMsg gatewayv1.SingleMessage
	if err := protojson.Unmarshal(rawMsg.Data, &singleMsg); err != nil {
		log.Printf("解析 SingleMessage 失败: %v", err)
		return
	}
}
