package cmd

import (
	"fmt"
	"github.com/labbcb/brave/client"
	"github.com/labbcb/brave/search"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func init() {
	searchCmd.Flags().String("host", "http://localhost:8080", "URL to BraVE server.")
	viper.BindPFlag("host", searchCmd.Flags().Lookup("host"))

	searchCmd.Flags().String("dataset", "", "Dataset name.")

	searchCmd.Flags().String("reference", "", "Genome version.")

	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for variants using queries",
	Long: `BraVE supports multiple types of queries:
	Gene symbol (SCN1A) returns variants that were annotated with a matching gene name.
	Genomic range (1:15000-16000) returns variants that are inside the range (1-based, half-open).
	Genomic position (1:12345) returns a single variant that have the same position (1-based)
	dbSNP ID (rs12345) returns a single variant that were annotated with this identifier.`,
	Run: func(cmd *cobra.Command, args []string) {
		var qs []*search.Query
		for _, text := range args {
			q := search.Parse(text)
			q.DatasetID = viper.GetString("dataset")
			q.AssemblyID = viper.GetString("assembly")
			qs = append(qs, q)
		}

		c := &client.Client{Host: viper.GetString("host")}
		resp, err := c.SearchVariants(&search.Input{Queries: qs})
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range resp.Variants {
			fmt.Println(v)
		}
	},
}
