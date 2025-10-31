package main

import (
	"os"
	toolerCmd "rain-im-server/internal/tooler/cmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tooler",
	Short: "A brief description of your application",
	Long:  `etcd get --key "key1" -a 8.148.84.185:2379`,
}

func init() {
	rootCmd.AddCommand(toolerCmd.EtcdCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
