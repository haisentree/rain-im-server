package service

import (
	"context"
	"fmt"
	"rain-im-server/global"
	gatewayv1 "rain-im-server/protogo/gateway/v1"

	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/encoding/protojson"
)

type MessageHandle struct {
	Validater *validator.Validate

	sender MessageSender // 持有消息发送器接口
}

func NewMessageHandle(sender MessageSender) *MessageHandle {
	return &MessageHandle{
		Validater: validator.New(),
		sender:    sender,
	}
}

func (m *MessageHandle) SingleMessageHandle(rawMsg *gatewayv1.RawMessage) {
	ctx := context.Background()

	var singleMsg gatewayv1.SingleMessage
	if err := protojson.Unmarshal(rawMsg.Data, &singleMsg); err != nil {
		fmt.Println("解析 SingleMessage 失败:", err)
		return
	}

	// 获取seq
	seqKey := fmt.Sprintf("%s%s-%s", global.ClientConversationStringKey, singleMsg.SourceId.ToUUID().String(), singleMsg.TargetId.ToUUID().String())

	// 先检查 key 是否存在
	// exists, err := global.Redis.Exists(ctx, seqKey).Result()
	// if err != nil {
	// 	fmt.Println("检查key是否存在失败:", err)
	// 	return
	// }

	// if exists == 0 {
	// 	// key 不存在，这是首次创建
	// 	fmt.Println("首次创建会话，seq从1开始")
	// }

	seq, err := global.Redis.Incr(ctx, seqKey).Uint64()
	if err != nil {
		fmt.Println("获取消息序号失败:", err)
		return
	}

	singleMsg.Seq = seq

	fmt.Println("消息序号已更新:", seq)

	rawMsgByte, err := protojson.Marshal(&singleMsg)
	if err != nil {
		fmt.Println("序列化 rawMsg 失败:", err)
		return
	}

	// 1.写入到持久化存储队列
	rawMsg.Type = gatewayv1.Message_MESSAGE_DB_SAVE
	pubSaveTheme := fmt.Sprintf(global.GatewayMessageSaveTheme, singleMsg.TargetId.ToUUID().String()[0:2]) // 取前两位作为分片
	if ack, err := global.NatsJS.Publish(pubSaveTheme, rawMsgByte); err != nil {
		fmt.Println("发布消息到 NATS 失败:", err)
		return
	} else {
		fmt.Printf("消息已发布到 NATS: %v\n", ack)
	}
	fmt.Println("消息已发布到 NATS 主题:", pubSaveTheme)

	// 2.消息会话seq+1
	pubSeqIncreaseTheme := fmt.Sprintf(global.GatewayMessageSeqIncreaseTheme, singleMsg.TargetId.ToUUID().String()[0:2]) // 取前两位作为分片
	if ack, err := global.NatsJS.Publish(pubSeqIncreaseTheme, rawMsgByte); err != nil {
		fmt.Println("发布消息到 NATS 失败:", err)
		return
	} else {
		fmt.Printf("消息已发布到 NATS: %v\n", ack)
	}
	fmt.Println("消息已发布到 NATS 主题:", pubSeqIncreaseTheme)

	// 3.本地和远程连接消息发送
	m.sender.WriteToLocalClient(singleMsg.TargetId.ToUUID().String(), rawMsg)
	rawMsg.Type = gatewayv1.Message_MESSAGE_SINGLE_GATEWAY_RELAY // 改变消息类型再发送
	m.sender.WriteToRemoteClient(singleMsg.TargetId.ToUUID().String(), rawMsg)
}

func (m *MessageHandle) RelayGatewaySingleMessageHandle(rawMsg *gatewayv1.RawMessage, rw *ResponseWriter) {
	var relaySingleMsg gatewayv1.SingleMessage
	if err := protojson.Unmarshal(rawMsg.Data, &relaySingleMsg); err != nil {
		fmt.Println("解析 relaySingleMsg 失败:", err)
		return
	}

	rawMsg.Type = gatewayv1.Message_MESSAGE_SINGLE // 还原消息类型再发送
	rw.WriteToLocalClient(relaySingleMsg.TargetId.ToUUID().String(), rawMsg)
}

func GroupHandleMessage(rawMsg *gatewayv1.RawMessage) {
	var groupMsg gatewayv1.GroupMessage
	if err := protojson.Unmarshal(rawMsg.Data, &groupMsg); err != nil {
		fmt.Println("解析 groupMsg 失败:", err)
		return
	}

	// groupMsg.Seq = 1

}
