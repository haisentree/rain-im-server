package service

import (
	"fmt"
	"log"
	"rain-im-server/global"
	gatewayv1 "rain-im-server/protogo/gateway/v1"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
)

// ResponseWriter 封装消息发送逻辑
type ResponseWriter struct {
	connManager *ConnectionManager // 连接管理器，用于给其他连接发消息
	// TOOD: GatewayServer、ResponseWriter依赖了ConnectionManager
	// ResponseWriter只使用了connManager的部分方法,应该抽象出一个接口
}

// NewResponseWriter 创建 ResponseWriter 实例
func NewResponseWriter(cm *ConnectionManager) *ResponseWriter {
	return &ResponseWriter{
		connManager: cm,
	}
}

// WriteToAllUserDevices 发送消息给指定用户的所有连接（多端同步）
func (rw *ResponseWriter) WriteToLocalClient(clientId string, msg *gatewayv1.RawMessage) error {
	conns := rw.connManager.GetClientConns(clientId)
	if len(conns) == 0 {
		return nil
	}

	// 客户端连接存在当前网关
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

func (rw *ResponseWriter) WriteToRemoteClinet(clientId string, msg *gatewayv1.RawMessage) {
	connDetails, err := rw.connManager.GetRemoteConnsInfo(clientId)
	if err != nil {
		log.Printf("get remote connections for client %s failed: %v", clientId, err)
		return
	}
	if len(connDetails) == 0 {
		return
	}

	// 发送消息给远程连接
	for _, detail := range connDetails {
		// 实现发送逻辑
		fmt.Println("123", detail)
		// 只转发单发消息,转发时候类型改成转发消息
		targetTopic := fmt.Sprintf(global.GatewayRelayMessageTheme, detail.GatewayUUID)

		err := global.Nats.Publish(targetTopic, msg.Data)
		if err != nil {
			log.Printf("publish relay message to topic %s failed: %v", targetTopic, err)
			continue
		}
	}

	// 	type ConnDetail struct {
	// 	ClientId    string `json:"clientId"`
	// 	PlatformId  string `json:"platformId"`
	// 	CountKey    string `json:"countKey"`
	// 	GatewayUUID string `json:"gatewayUUID"`
	// 	CreatedAt   int64  `json:"createdAt"`
	// }

}

func (rw *ResponseWriter) pubToNats(client, gateway string, relayMsg *gatewayv1.RelayMessage) error {
	return nil
}

// // WriteRaw 发送原始字节消息（需自行序列化）
// func (rw *ResponseWriter) WriteRaw(conn *WSClient, msgType int, data []byte) error {
// 	return conn.WriteMessage(msgType, data)
// }
