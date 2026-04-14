package service

import (
	"encoding/json"
	"fmt"
	"rain-im-server/global"
	gatewayv1 "rain-im-server/protogo/gateway/v1"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"google.golang.org/protobuf/encoding/protojson"
)

type MessageHandle struct {
	Decoder   *schema.Decoder
	Validater *validator.Validate
}

// ResponseWriter不作为参数传递,而是MessageHandle持有ResponseWriter的引用,这样MessageHandle就可以调用ResponseWriter的方法来发送消息
// 抽象出接口,让MessageHandle依赖接口而不是具体的ResponseWriter实现,这样可以降低耦合度,提高代码的可测试性和可维护性
func NewMessageHandle() *MessageHandle {
	decoder := schema.NewDecoder()
	// decoder.IgnoreUnknownKeys(true)
	// decoder.SetAliasTag("schema")
	return &MessageHandle{
		Decoder:   decoder,
		Validater: validator.New(),
	}
}

func (m *MessageHandle) SingleMessageHandle(rawMsg *gatewayv1.RawMessage, rw *ResponseWriter) {
	var singleMsg gatewayv1.SingleMessage
	if err := protojson.Unmarshal(rawMsg.Data, &singleMsg); err != nil {
		fmt.Println("解析 SingleMessage 失败:", err)
		return
	}

	fmt.Println(singleMsg.SourceId)
	fmt.Println(singleMsg.TargetId)
	fmt.Println(singleMsg.Content)

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
	rw.WriteToLocalClient(singleMsg.TargetId.ToUUID().String(), rawMsg)
	rawMsg.Type = gatewayv1.Message_MESSAGE_RELAY_GATEWAY_SINGLE // 改变消息类型再发送
	rw.WriteToRemoteClinet(singleMsg.TargetId.ToUUID().String(), rawMsg)
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

// 转发消息的内容是单方消息结构,只发送到本地的连接
func (m *MessageHandle) RelayMessageHandle(data []byte, conn *WSClient, rw *ResponseWriter) {
	var relayMsg gatewayv1.RelayMessage
	if err := protojson.Unmarshal(data, &relayMsg); err != nil {
		fmt.Println("解析 RelayMessage 失败:", err)
		return
	}

	fmt.Println(relayMsg.SourceId)
	fmt.Println(relayMsg.TargetId)
	fmt.Println(relayMsg.Content)

	// rawMsg := &gatewayv1.RawMessage{
	// 	Type: gatewayv1.Message_MESSAGE_CONFIRM,
	// 	Data: data,
	// }

	// rw.WriteToClient(relayMsg.TargetId.ToUUID().String(), false, rawMsg)
}
