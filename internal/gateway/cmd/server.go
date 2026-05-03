package cmd

import (
	"rain-im-server/global"
	"rain-im-server/internal/gateway/service"

	"github.com/spf13/cobra"
)

const GatewayName = "gateway-1"

func ServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "run rain-im-server",
		Run: func(cmd *cobra.Command, args []string) {
			RunServer()
		},
	}
}

func RunServer() {
	gatewayServer, err := service.NewGatewayServer("0.0.0.0:5173")
	if err != nil {
		panic("RunServer err")
	}
	defer global.CloseService()

	// ctx := context.Background()
	// global.Redis.Set(ctx, "gateway:name", "123", 5*time.Minute)
	// key := global.Redis.Get(ctx, "gateway:name")
	// fmt.Println("key:", key.Val())

	gatewayServer.Run()
}
