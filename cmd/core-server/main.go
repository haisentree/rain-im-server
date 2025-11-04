package main

import (
	"os"
	coreCmd "rain-im-server/internal/core/cmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tooler",
	Short: "A brief description of your application",
	Long:  `etcd get --key "key1" -a 8`,
}

func init() {
	rootCmd.AddCommand(coreCmd.ServerCmd())
	rootCmd.AddCommand(coreCmd.VersionCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
