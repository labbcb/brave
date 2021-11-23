package cmd

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/labbcb/brave/client"
	"github.com/labbcb/brave/variant"

	"github.com/labbcb/brave/vcf"
	"github.com/spf13/cobra"
)

var dontFilter, dryRun bool

func init() {
	importCmd.Flags().StringVar(&host, "host", "http://localhost:8080", "URL to BraVE server.")

	importCmd.Flags().StringVar(&datasetID, "dataset", "", "Dataset name.")
	importCmd.MarkFlagRequired("dataset")

	importCmd.Flags().StringVar(&assemblyID, "assembly", "", "Genome version.")
	importCmd.MarkFlagRequired("assembly")

	importCmd.Flags().StringVar(&username, "username", "admin", "User name.")

	importCmd.Flags().StringVar(&password, "password", "", "Password.")

	importCmd.Flags().BoolVar(&dontFilter, "dont-filter", false, "Don't filter variants by FILTER column.")
	importCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Just check VCF without connecting to server.")

	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import genomic variants from VCF files to database",
	Long: `BraVE importer supports variant data in Variant Call Format (VCF) files.
	The server should be running, see brave help server.
	Existing variants that have the same dataset and reference genome are not removed by default.
	If some variant in VCF file has the same Reference Name (chromosome) and Position then it will panic.
	See brave help remove to delete previous data before importing.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			if err := importVcf(file); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func isGzip(file string) (bool, error) {
	f, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer f.Close()

	br := bufio.NewReader(f)
	b, err := br.Peek(2)
	if err != nil {
		return false, err
	}

	return b[0] == 31 && b[1] == 139, nil
}

func importVcf(file string) error {
	var r io.Reader
	r, err := os.Open(file)
	if err != nil {
		return err
	}

	gz, err := isGzip(file)
	if err != nil {
		return err
	}
	if gz {
		r, err = gzip.NewReader(r)
		if err != nil {
			return err
		}
	}

	c := &client.Client{
		Host:     host,
		Username: username,
		Password: password,
	}

	importVariant := func(v *variant.Variant) error {
		return c.InsertVariant(v)
	}

	if dryRun {
		importVariant = func(v *variant.Variant) error { return nil }
	}

	doFilter := !dontFilter
	summary, err := vcf.IterateOver(r, doFilter, datasetID, assemblyID, importVariant)

	fmt.Println("Total variants:", summary.TotalVariants)
	if doFilter {
		fmt.Println("Passed variants:", summary.PassedVariants)
	}

	return err
}
