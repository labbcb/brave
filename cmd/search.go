package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/labbcb/brave/client"
	"github.com/labbcb/brave/search"
	"github.com/labbcb/brave/variant"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
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
		case "csv":
			if len(resp.Variants) == 0 {
				return
			}

			w := csv.NewWriter(os.Stdout)

			h := []string{
				"dataset",
				"assembly",
				"ns",
				"total",
				"chrom",
				"pos",
				"id",
				"ref",
				"alt",
				"af",
				"dp",
				"gq",
				"gene",
			}
			if err := w.Write(h); err != nil {
				log.Fatal(err)
			}

			for _, v := range resp.Variants {
				s := []string{
					v.DatasetID,
					v.AssemblyID,
					strconv.Itoa(v.SampleCount),
					strconv.Itoa(int(v.TotalSamples)),
					v.ReferenceName,
					strconv.Itoa(int(v.Start)),
					strings.Join(v.SnpIds, ";"),
					v.ReferenceBases,
					strings.Join(v.AlternateBases, ";"),
					joinFloats(v.AlleleFrequency, ";"),
					joinDistribution(v.Coverage),
					joinDistribution(v.GenotypeQuality),
					strings.Join(v.GeneSymbol, ";"),
				}
				if err := w.Write(s); err != nil {
					log.Fatal(err)
				}
			}
			w.Flush()
		default:
			for _, v := range resp.Variants {
				fmt.Println(v)
			}
		}
	},
}

func joinFloats(fs []float32, sep string) string {
	var a []string
	for _, f := range fs {
		a = append(a, fmt.Sprintf("%f", f))
	}
	return strings.Join(a, sep)
}

func joinDistribution(d *variant.Distribution) string {
	return fmt.Sprintf("%f;%f;%f;%f;%f;%f",
		d.Min,
		d.Q25,
		d.Median,
		d.Q75,
		d.Max,
		d.Mean)
}
