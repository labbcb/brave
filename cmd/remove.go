package cmd

import (
	"github.com/labbcb/brave/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	removeCmd.Flags().String("host", "http://localhost:8080", "URL to BraVE server.")
	viper.BindPFlag("host", removeCmd.Flags().Lookup("host"))

	removeCmd.Flags().String("dataset", "", "Dataset name.")
	removeCmd.MarkFlagRequired("dataset")

	removeCmd.Flags().String("assembly", "", "Genome version.")
	removeCmd.MarkFlagRequired("assembly")

	removeCmd.Flags().String("username", "admin", "User name.")
	viper.BindPFlag("username", removeCmd.Flags().Lookup("username"))

	removeCmd.Flags().String("password", "", "Password.")
	viper.BindPFlag("password", removeCmd.Flags().Lookup("password"))

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
			Host:     viper.GetString("host"),
			Username: viper.GetString("username"),
			Password: viper.GetString("password"),
		}

		datasetID := viper.GetString("dataset")
		assemblyID := viper.GetString("assembly")
		if err := c.RemoveVariants(datasetID, assemblyID); err != nil {
			log.Fatal(err)
		}
	},
}
