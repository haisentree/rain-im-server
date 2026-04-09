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

func NewMessageHandle() *MessageHandle {
	decoder := schema.NewDecoder()
	// decoder.IgnoreUnknownKeys(true)
	// decoder.SetAliasTag("schema")
	return &MessageHandle{
		Decoder:   decoder,
		Validater: validator.New(),
	}
}

func (m *MessageHandle) SingleMessageHandle(data []byte, conn *WSClient, rw *ResponseWriter) {
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

	rw.WriteToClient(singleMsg.TargetId.ToUUID().String(), rawMsg)
	// 从server中获取conn
	// 写入内容
}
