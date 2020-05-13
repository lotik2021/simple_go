package cmd

import (
	"bitbucket.movista.ru/maas/maasapi/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "maas",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log.Info("YO!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
