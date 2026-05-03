package cmd

import (
	"rain-im-server/internal/core/worker"

	"github.com/spf13/cobra"
)

func WorkerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "worker",
		Short: "run manage worker",
		Run: func(cmd *cobra.Command, args []string) {
			RunWorker()
		},
	}
}

func RunWorker() {
	worker := worker.NewWorker()
	worker.Run()
}
