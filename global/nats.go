package global

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

const (
	// 客户端连接不在网关服务器上，就投放在消息队列中，由消息队列分发给网关服务器
	GatewayRelayMessageTheme = "gateway.message.relay.%s"
)

const (
	GatewayStreamNameMessage       = "GATEWAY-PERSIST-MESSAGE"
	GatewayMessageSeqIncreaseTheme = "gateway.message.seq.increase.%s"
	GatewayMessageSeqIncreaseQueue = "gateway-message-seq-increase" //  这个队列中单个消费者顺序消费
	GatewayMessageSaveTheme        = "gateway.message.save.%s"
	GatewayMessageSaveQueue        = "gateway-message-save"
)

func initJetStream(js nats.JetStreamContext) error {
	subjects := []string{
		fmt.Sprintf(GatewayRelayMessageTheme, "*"), // 匹配任意网关ID
		fmt.Sprintf(GatewayMessageSeqIncreaseTheme, "*"),
		fmt.Sprintf(GatewayMessageSaveTheme, "*"),
	}

	_, err := js.StreamInfo(GatewayStreamNameMessage)
	if err == nil {
		return err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      GatewayStreamNameMessage,
		Subjects:  subjects,
		Storage:   nats.FileStorage,
		Retention: nats.WorkQueuePolicy,
	})
	if err != nil {
		log.Fatalf("创建流失败: %v", err)
	}

	return err
}
