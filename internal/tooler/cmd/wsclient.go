package tooler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	gatewayv1 "rain-im-server/proto/protogo/gateway/v1"
	basev1 "rain-im-server/protogo/base/v1"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// WSClientCmd 主命令（和你的 etcdCmd 一模一样风格）
func WSClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ws",
		Short: "websocket client tool",
	}

	cmd.AddCommand(WSClientSendCmd())
	return cmd
}

// 你的消息结构体（完全按你写的格式）
type WSRequest struct {
	Type string        `json:"type"`
	Data WSRequestData `json:"data"`
}

type WSRequestData struct {
	SourceId string `json:"source_id"`
	TargetId string `json:"target_id"`
	Content  string `json:"content"`
}

// WSClientSendCmd 发送 + 持续接收
func WSClientSendCmd() *cobra.Command {
	var address string

	cmd := &cobra.Command{
		Use:   "send",
		Short: "send websocket message & keep receiving",
		Run: func(cmd *cobra.Command, args []string) {
			// 1. 连接 WS
			conn, _, err := websocket.DefaultDialer.DialContext(context.Background(), address, nil)
			if err != nil {
				log.Fatal("ws connect error:", err)
			}
			defer conn.Close()

			log.Println("✅ websocket connected:", address)

			// ======================
			// 核心：协程 1 → 持续接收服务器消息
			// ======================
			go func() {
				for {
					// 阻塞等待接收
					msgType, msgBody, err := conn.ReadMessage()
					if err != nil {
						log.Println("\n🔌 server disconnected:", err)
						return
					}

					// 打印收到的消息
					if msgType == websocket.TextMessage {
						log.Printf("\n📥 received server message:\n%s\n", string(msgBody))
						log.Print("> enter json to send: ")
					}
				}
			}()

			// ======================
			// 主线程 → 持续读取用户输入发送
			// ======================
			log.Println("\n> enter json to send (type 'exit' to quit):")
			scanner := bufio.NewScanner(os.Stdin)

			var jsonBuilder strings.Builder

			for {
				print("> ")
				if !scanner.Scan() {
					break
				}

				line := scanner.Text()
				trimmed := strings.TrimSpace(line)

				// 退出
				if trimmed == "exit" {
					log.Println("👋 exit...")
					break
				}

				if trimmed == "" {
					fullJSON := jsonBuilder.String()
					if fullJSON == "" {
						continue
					}

					var msg WSRequest
					if err := json.Unmarshal([]byte(fullJSON), &msg); err != nil {
						log.Println("❌ json error:", err)
						jsonBuilder.Reset()
						continue
					}

					compressed, _ := json.Marshal(msg)
					log.Println("📤 待发送内容：", string(compressed))

					sourceUUID := basev1.NewUUID()
					sourceUUID.FromString(msg.Data.SourceId)
					targetUUID := basev1.NewUUID()
					targetUUID.FromString(msg.Data.TargetId)

					var singleMsg gatewayv1.SingleMessage

					singleMsg.SourceId = sourceUUID
					singleMsg.TargetId = targetUUID
					singleMsg.Content = msg.Data.Content

					singleBytes, err := protojson.Marshal(&singleMsg)
					if err != nil {
						log.Println("error1:", err.Error())
						continue
					}
					rawMsg := &gatewayv1.RawMessage{
						Type: gatewayv1.Message(gatewayv1.Message_value[msg.Type]),
						Data: singleBytes,
					}
					fmt.Println("rawMsg:", rawMsg)

					// 序列化成最终发送JSON（自动压缩）
					sendBytes, err := protojson.Marshal(rawMsg)
					if err != nil {
						log.Println("❌ sendBytes error:", err)
						continue
					}

					err = conn.WriteMessage(websocket.TextMessage, sendBytes)
					if err != nil {
						log.Println("❌ send error:", err)
						continue
					} else {
						log.Println("✅ send success!")
					}
					jsonBuilder.Reset()
					continue
				}

				jsonBuilder.WriteString(line)
				jsonBuilder.WriteByte(' ')
			}

			if err := scanner.Err(); err != nil {
				log.Fatal("❌ input error:", err)
			}
		},
	}

	// 命令行参数
	cmd.Flags().StringVarP(&address, "address", "a", "", "websocket server addr (ws://ip:port/path)")
	_ = cmd.MarkFlagRequired("address")

	return cmd
}
