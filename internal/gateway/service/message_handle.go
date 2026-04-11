package service

import (
	"fmt"
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

func (m *MessageHandle) SingleMessageHandle(data []byte, rw *ResponseWriter) {
	var singleMsg gatewayv1.SingleMessage
	if err := protojson.Unmarshal(data, &singleMsg); err != nil {
		fmt.Println("解析 SingleMessage 失败:", err)
		return
	}

	fmt.Println(singleMsg.SourceId)
	fmt.Println(singleMsg.TargetId)
	fmt.Println(singleMsg.Content)

	rawMsg := &gatewayv1.RawMessage{
		Type: gatewayv1.Message_MESSAGE_CONFIRM,
		Data: data,
	}

	// conn.WriteToSelf(rawMsg)

	// 1.写入到持久化存储队列
	// 2.尝试发送到本地连接,不存在则发送对应在线网关订阅的主题

	switch singleMsg.Method {
	case gatewayv1.Method_METHOD_ALL:
		rw.WriteToLocalClient(singleMsg.TargetId.ToUUID().String(), rawMsg)
		rw.WriteToRemoteClinet(singleMsg.TargetId.ToUUID().String(), rawMsg)
	case gatewayv1.Method_METHOD_LOCAL:
		rw.WriteToRemoteClinet(singleMsg.TargetId.ToUUID().String(), rawMsg)
	default:
		panic("err")
	}

	// 从server中获取conn
	// 写入内容
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
