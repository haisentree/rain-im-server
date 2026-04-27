package service

import (
	"encoding/json"
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
	var singleMsg gatewayv1.SingleMessage
	if err := protojson.Unmarshal(rawMsg.Data, &singleMsg); err != nil {
		fmt.Println("解析 SingleMessage 失败:", err)
		return
	}

	rawMsgByte, err := json.Marshal(rawMsg)
	if err != nil {
		fmt.Println("序列化 rawMsg 失败:", err)
		return
	}
	// 1.写入到持久化存储队列
	pubSaveTheme := fmt.Sprintf(global.GatewayMessageSaveTheme, singleMsg.TargetId.ToUUID().String()[0:2]) // 取前两位作为分片
	if ack, err := global.NatsJS.Publish(pubSaveTheme, rawMsgByte); err != nil {
		fmt.Println("发布消息到 NATS 失败:", err)
		return
	} else {
		fmt.Printf("消息已发布到 NATS: %v\n", ack)
	}

	// 2.消息会话seq+1
	pubSeqIncreaseTheme := fmt.Sprintf(global.GatewayMessageSeqIncreaseTheme, singleMsg.TargetId.ToUUID().String()[0:2]) // 取前两位作为分片
	if ack, err := global.NatsJS.Publish(pubSeqIncreaseTheme, rawMsgByte); err != nil {
		fmt.Println("发布消息到 NATS 失败:", err)
		return
	} else {
		fmt.Printf("消息已发布到 NATS: %v\n", ack)
	}

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
