package cmd

import (
	"rain-im-server/internal/core/service"

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
	httpServer := service.Init()
	httpServer.Addr = "0.0.0.0:8081"

	httpServer.ListenAndServe()
}
