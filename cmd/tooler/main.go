package main

import (
	"os"
	toolerCmd "rain-im-server/internal/tooler/cmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tooler",
	Short: "A brief description of your application",
	Long:  `etcd`,
}

// etcd get --key "key1" -a 8.148.84.185:2379
// etcd getPrefix --key "/" -a 8.148.84.185:2379
// etcd put --key "/key" -v "value2" -a 8.148.84.185:2379

// ws send -a 'ws://8.148.84.185:5173/gateway?token={"client_id":"111","platform":2}'

func init() {
	rootCmd.AddCommand(toolerCmd.EtcdCmd())
	rootCmd.AddCommand(toolerCmd.WSClientCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
