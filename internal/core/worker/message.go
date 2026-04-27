package worker

import (
	"log"
	"rain-im-server/global"

	"github.com/nats-io/nats.go"
)

type MessageWorker struct {
}

func NewMessageWorker() *MessageWorker {
	return &MessageWorker{}
}

func (w *MessageWorker) Run() {
}

func (w *MessageWorker) Subscribe() {
	sub2, err := global.NatsJS.QueueSubscribe("orders-2.>", "worker-group", func(msg *nats.Msg) {
		log.Printf("收到订单: %s", string(msg.Data))
		msg.Ack()
	}, nats.ManualAck())
	if err != nil {
		log.Fatal(err)
	}
	defer sub2.Unsubscribe()
}
