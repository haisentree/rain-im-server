package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "1.0.0"

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "run manage server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version:%s", VERSION)
		},
	}
}
