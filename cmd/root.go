package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "brave",
	Short: "BraVE - BIPMed Variant Explorer",
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("brave")
}

// Execute starts command line parser.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
