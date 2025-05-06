package main

import (
	cpustress "github.com/QQGoblin/StressMaker/cmd/cpu"
	"github.com/QQGoblin/StressMaker/cmd/excel"
	"github.com/QQGoblin/StressMaker/cmd/freq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stress",
	Short: "Stress Maker",
}

func init() {
	rootCmd.AddCommand(cpustress.Command)
	rootCmd.AddCommand(freq.Command)
	rootCmd.AddCommand(excel.Command)
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("exit!")
	}
}
