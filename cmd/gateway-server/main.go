package main

import (
	"os"
	gatewayCmd "rain-im-server/internal/gateway/cmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tooler",
	Short: "A brief description of your application",
	Long:  `etcd get --key "key1" -a 8`,
}

func init() {
	rootCmd.AddCommand(gatewayCmd.ServerCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
