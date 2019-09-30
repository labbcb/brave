package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var datasetID, assemblyID, database, host, address, username, password string

var rootCmd = &cobra.Command{
	Use:   "brave",
	Short: "BraVE - BIPMed Variant Explorer",
}

// Execute starts command line parser.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
