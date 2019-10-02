package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/labbcb/brave/client"
	"github.com/labbcb/brave/search"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var format string

func init() {
	searchCmd.Flags().StringVar(&host, "host", "http://localhost:8080", "URL to BraVE server.")
	searchCmd.Flags().StringVar(&datasetID, "dataset", "", "Dataset name.")
	searchCmd.Flags().StringVar(&assemblyID, "assembly", "", "Genome version.")
	searchCmd.Flags().StringVar(&format, "format", "console", "Output format.")

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

		c := &client.Client{Host: host}
		resp, err := c.SearchVariants(&search.Input{Queries: qs})
		if err != nil {
			log.Fatal(err)
		}

		switch format {
		case "json":
			if err := json.NewEncoder(os.Stdout).Encode(resp.Variants); err != nil {
				log.Fatal(err)
			}
		default:
			for _, v := range resp.Variants {
				fmt.Println(v)
			}
		}
	},
}
