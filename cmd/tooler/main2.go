package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// 连接到 NATS 服务器（默认地址 nats://localhost:4222）
	nc, err := nats.Connect("nats://root:haisen123@aliyun.haisentree.top:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// 可选：设置连接事件回调
	nc.SetReconnectHandler(func(conn *nats.Conn) {
		log.Println("重新连接成功")
	})

	// 订阅主题
	sub, err := nc.Subscribe("foo", func(msg *nats.Msg) {
		log.Printf("收到消息: %s", string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// 发布消息
	err = nc.Publish("foo", []byte("Hello NATS!"))
	if err != nil {
		log.Fatal(err)
	}

	err = nc.Publish("foo", []byte("123"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("123")

	select {}
}
