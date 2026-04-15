package service

import (
	"fmt"
	"log"
	"rain-im-server/global"
	gatewayv1 "rain-im-server/protogo/gateway/v1"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
)

type MessageSender interface {
	// WriteToLocalClient 发送消息给本地网关上的指定用户所有连接
	WriteToLocalClient(clientId string, msg *gatewayv1.RawMessage) error
	// WriteToRemoteClient 发送消息给远程网关上的指定用户所有连接
	WriteToRemoteClient(clientId string, msg *gatewayv1.RawMessage)
}

// ResponseWriter 封装消息发送逻辑
type ResponseWriter struct {
	connGetter ConnectionGetter
}

func NewResponseWriter(cm ConnectionGetter) *ResponseWriter {
	return &ResponseWriter{
		connGetter: cm,
	}
}

// WriteToAllUserDevices 发送消息给指定用户的所有连接（多端同步）
func (rw *ResponseWriter) WriteToLocalClient(clientId string, msg *gatewayv1.RawMessage) error {
	conns := rw.connGetter.GetLocalClientConns(clientId)
	if len(conns) == 0 {
		return nil
	}

	data, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	var lastErr error
	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			lastErr = err
			log.Println("err:", err)
			continue
		}
	}
	return lastErr
}

func (rw *ResponseWriter) WriteToRemoteClient(clientId string, msg *gatewayv1.RawMessage) {
	connDetails, err := rw.connGetter.GetRemoteConnsInfo(clientId)
	if err != nil {
		log.Printf("get remote connections for client %s failed: %v", clientId, err)
		return
	}

	for _, detail := range connDetails {
		targetTopic := fmt.Sprintf(global.GatewayRelayMessageTheme, detail.GatewayUUID)

		err := global.Nats.Publish(targetTopic, msg.Data)
		if err != nil {
			log.Printf("publish relay message to topic %s failed: %v", targetTopic, err)
			continue
		}
	}
}

// func (rw *ResponseWriter) pubToNats(client, gateway string, relayMsg *gatewayv1.RelayMessage) error {
// 	return nil
// }
