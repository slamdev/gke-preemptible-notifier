package internal

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const rootCmdName = "gke-preemptible-notifier"

var rootCmd = &cobra.Command{
	Use:   rootCmdName,
	Short: rootCmdName + " Sends notifications when preemptible node got killed",
	Long: rootCmdName + ` Sends notifications when preemptible node got killed

 Find more information at: https://github.com/slamdev/` + rootCmdName,
}

func init() {
	rootCmd.PersistentFlags().String("log-format", "text", "log format (json or text)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn or error)")
	cobra.OnInitialize(func() {
		fillWithEnvVars(rootCmd.Flags())
	}, setup)
}

func setup() {
	lvl, err := rootCmd.Flags().GetString("log-level")
	if err != nil {
		logrus.WithError(err).Fatalf("failed to get %v flag", "log-level")
	}
	format, err := rootCmd.Flags().GetString("log-format")
	if err != nil {
		logrus.WithError(err).Fatalf("failed to get %v flag", "log-format")
	}
	if err := initLogger(lvl, format); err != nil {
		logrus.WithError(err).Fatal("failed init logger")
	}
}

func initLogger(lvl string, format string) error {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return fmt.Errorf("failed to parse %v into log level. %w", lvl, err)
	}
	logrus.SetLevel(level)

	if format == "text" {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else if format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		return fmt.Errorf("unsupported log format %v. %w", format, err)
	}

	return nil
}

func ExecuteCmd() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("failed to run application")
	}
}
