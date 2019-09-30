package cmd

import (
	"fmt"
	"github.com/labbcb/brave/client"
	"github.com/labbcb/brave/search"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	searchCmd.Flags().StringVar(&host, "host", "http://localhost:8080", "URL to BraVE server.")
	searchCmd.Flags().StringVar(&datasetID, "dataset", "bipmed", "Dataset name")
	searchCmd.Flags().StringVar(&assemblyID, "reference", "GRCh38", "Genome version")
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
			q.DatasetID = datasetID
			q.AssemblyID = assemblyID
			qs = append(qs, q)
		}

		c := &client.Client{
			Host:     host,
			Username: username,
			Password: password,
		}
		resp, err := c.SearchVariants(&search.Input{Queries: qs})
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range resp.Variants {
			fmt.Println(v)
		}
	},
}
