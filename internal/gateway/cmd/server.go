package cmd

import (
	_ "rain-im-server/internal/gateway/global"
	"rain-im-server/internal/gateway/service"

	"github.com/spf13/cobra"
)

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
	gatewayServer := service.NewGatewayServer()
	gatewayServer.Addr = ":5173"

	gatewayServer.Run()
}
