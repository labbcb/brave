package cmd

import (
	"github.com/labbcb/brave/client"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	removeCmd.Flags().StringVar(&host, "host", "http://localhost:8080", "URL to BraVE server.")
	removeCmd.Flags().StringVar(&datasetID, "dataset", "bipmed", "Dataset name")
	removeCmd.Flags().StringVar(&assemblyID, "reference", "GRCh38", "Genome version")
	removeCmd.Flags().StringVar(&username, "username", "admin", "User name.")
	removeCmd.Flags().StringVar(&password, "password", "", "Password.")
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Delete variants from database given a dataset and genome version",
	Long: `BraVE will delete all genomics variants that have have the same dataset AND genome version.
	If genome version is not specified then it will removes variants that matches dataset and vice-versa.
	If none of them are specified, which is default, remove all variants.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := &client.Client{
			Host:     host,
			Username: username,
			Password: password,
		}
		if err := c.RemoveVariants(datasetID, assemblyID); err != nil {
			log.Fatal(err)
		}
	},
}
