package cmd

import (
	"bufio"
	"compress/gzip"
	"github.com/labbcb/brave/client"
	"github.com/labbcb/brave/variant"
	"io"
	"log"
	"os"

	"github.com/labbcb/brave/vcf"
	"github.com/spf13/cobra"
)

func init() {
	importCmd.Flags().StringVar(&host, "host", "http://localhost:8080", "URL to BraVE server.")
	importCmd.Flags().StringVar(&datasetID, "dataset", "bipmed", "Dataset name")
	importCmd.Flags().StringVar(&assemblyID, "reference", "GRCh38", "Genome version")
	importCmd.Flags().StringVar(&username, "username", "admin", "User name.")
	importCmd.Flags().StringVar(&password, "password", "", "Password.")
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
	return vcf.IterateOver(r, datasetID, assemblyID, func(v *variant.Variant) error {
		return c.InsertVariant(v)
	})
}
