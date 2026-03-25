package service

import (
	"log"
	gatewayv1 "rain-im-server/protogo/gateway/v1"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
)

// ResponseWriter 封装消息发送逻辑
type ResponseWriter struct {
	connManager *ConnectionManager // 连接管理器，用于给其他连接发消息
}

// NewResponseWriter 创建 ResponseWriter 实例
func NewResponseWriter(cm *ConnectionManager) *ResponseWriter {
	return &ResponseWriter{
		connManager: cm,
	}
}

// WriteToAllUserDevices 发送消息给指定用户的所有连接（多端同步）
func (rw *ResponseWriter) WriteToAllUserDevices(clientId string, msg *gatewayv1.RawMessage) error {
	conns := rw.connManager.GetClientConns(clientId)
	if len(conns) == 0 {
		log.Println("client is not conn :", clientId)
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
			// 继续尝试发送给其他设备
		}
	}
	return lastErr
}

// // WriteRaw 发送原始字节消息（需自行序列化）
// func (rw *ResponseWriter) WriteRaw(conn *WSClient, msgType int, data []byte) error {
// 	return conn.WriteMessage(msgType, data)
// }
