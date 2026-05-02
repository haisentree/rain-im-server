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

type MessageWorker struct {
	subscriptions []*nats.Subscription
}

func NewMessageWorker() *MessageWorker {
	return &MessageWorker{}
}

func (w *MessageWorker) Run() {
	w.Subscribe()

	// 保持运行
	select {}
}

// 优雅关闭
func (w *MessageWorker) Stop() {
	for _, sub := range w.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			log.Printf("取消订阅失败: %v", err)
		}
	}
	log.Println("消息工作器已停止")
}

func (w *MessageWorker) Subscribe() {
	// 订阅所有分片的持久化主题
	// 使用通配符 * 匹配所有两位十六进制字符的分片（00-ff）
	saveSubject := fmt.Sprintf("%s.*", global.GatewayMessageSaveTheme[:len(global.GatewayMessageSaveTheme)-3])

	sub, err := global.NatsJS.QueueSubscribe(
		saveSubject,                    // 订阅主题，例如 "gateway.message.save.*"
		global.GatewayMessageSaveQueue, // 队列组名，多个实例负载均衡
		func(msg *nats.Msg) {
			w.handleMessage(msg)
		},
		nats.ManualAck(), // 手动确认
		nats.Durable(global.GatewayMessageSaveConsumer), // 持久化消费者
	)
	if err != nil {
		log.Fatalf("订阅消息持久化队列失败: %v", err)
	}

	w.subscriptions = append(w.subscriptions, sub)
	log.Printf("消息持久化工作器已启动，订阅主题: %s", saveSubject)

	// 如果需要同时订阅 Seq 增加队列，可以添加
	// seqSubject := fmt.Sprintf("%s.*", global.GatewayMessageSeqIncreaseTheme[:len(global.GatewayMessageSeqIncreaseTheme)-3])
	// seqSub, err := global.NatsJS.QueueSubscribe(
	// 	seqSubject,
	// 	"message-seq-worker-group",
	// 	func(msg *nats.Msg) {
	// 		w.handleSeqIncrease(msg)
	// 	},
	// 	nats.ManualAck(),
	// 	nats.Durable("message-seq-consumer"),
	// )
	// if err != nil {
	// 	log.Fatalf("订阅消息序号队列失败: %v", err)
	// }

	// w.subscriptions = append(w.subscriptions, seqSub)
	// log.Printf("消息序号工作器已启动，订阅主题: %s", seqSubject)

	// sub2, err := global.NatsJS.QueueSubscribe("orders-2.>", "worker-group", func(msg *nats.Msg) {
	// 	log.Printf("收到订单: %s", string(msg.Data))
	// 	msg.Ack()
	// }, nats.ManualAck())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer sub2.Unsubscribe()
}

// 处理消息持久化
func (w *MessageWorker) handleMessage(msg *nats.Msg) {
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

	// 3. 持久化到数据库
	if err := w.saveToDatabase(&singleMsg, &rawMsg); err != nil {
		log.Printf("保存消息到数据库失败: %v", err)
		// 这里可以选择 nack 重试
		// msg.Nak()
		return
	}

	// log.Printf("消息已持久化: MsgId=%s, TargetId=%s",
	// 	singleMsg.MsgId, singleMsg.TargetId.ToUUID().String())
}

// // 处理序号增加
// func (w *MessageWorker) handleSeqIncrease(msg *nats.Msg) {
// 	defer msg.Ack()

// 	var rawMsg gatewayv1.RawMessage
// 	if err := json.Unmarshal(msg.Data, &rawMsg); err != nil {
// 		log.Printf("解析 RawMessage 失败: %v", err)
// 		return
// 	}

// 	var singleMsg gatewayv1.SingleMessage
// 	if err := protojson.Unmarshal(rawMsg.Data, &singleMsg); err != nil {
// 		log.Printf("解析 SingleMessage 失败: %v", err)
// 		return
// 	}

// 	// 增加会话序号
// 	targetId := singleMsg.TargetId.ToUUID().String()
// 	if err := w.increaseSeq(targetId); err != nil {
// 		log.Printf("增加消息序号失败: targetId=%s, err=%v", targetId, err)
// 		return
// 	}

// 	log.Printf("消息序号已更新: TargetId=%s", targetId)
// }

// 保存消息到数据库
func (w *MessageWorker) saveToDatabase(singleMsg *gatewayv1.SingleMessage, rawMsg *gatewayv1.RawMessage) error {
	// TODO: 实现数据库保存逻辑
	// 示例：
	// db.Create(&model.Message{
	//     MsgId:    singleMsg.MsgId,
	//     FromId:   singleMsg.FromId.ToUUID().String(),
	//     TargetId: singleMsg.TargetId.ToUUID().String(),
	//     Content:  singleMsg.Content,
	//     SendTime: rawMsg.SendTime,
	// })

	return nil
}

// 增加消息序号
func (w *MessageWorker) increaseSeq(targetId string) error {
	// TODO: 实现序号自增逻辑
	// 示例：使用 Redis INCR
	// key := fmt.Sprintf("msg:seq:%s", targetId)
	// return redis.Incr(key).Err()

	return nil
}
