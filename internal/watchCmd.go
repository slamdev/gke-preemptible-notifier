package internal

import (
	"context"
	"github.com/spf13/cobra"
	"gke-preemptible-notifier/pkg/notifier"
)

var throwCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch preemptible node",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return notifier.Watch(context.Background(), cmd.OutOrStdout())
	},
}

func init() {
	cobra.OnInitialize(func() {
		fillWithEnvVars(throwCmd.Flags())
	})
	rootCmd.AddCommand(throwCmd)
}
