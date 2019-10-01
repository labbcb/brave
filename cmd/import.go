package cmd

import (
	"bufio"
	"compress/gzip"
	"github.com/labbcb/brave/client"
	"github.com/labbcb/brave/variant"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"

	"github.com/labbcb/brave/vcf"
	"github.com/spf13/cobra"
)

func init() {
	importCmd.Flags().String("host", "http://localhost:8080", "URL to BraVE server.")
	viper.BindPFlag("host", importCmd.Flags().Lookup("host"))

	importCmd.Flags().String("dataset", "", "Dataset name.")
	importCmd.MarkFlagRequired("dataset")

	importCmd.Flags().String("assembly", "", "Genome version.")
	importCmd.MarkFlagRequired("assembly")

	importCmd.Flags().String("username", "admin", "User name.")
	viper.BindPFlag("username", importCmd.Flags().Lookup("username"))

	importCmd.Flags().String("password", "", "Password.")
	viper.BindPFlag("password", importCmd.Flags().Lookup("password"))

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
		Host:     viper.GetString("host"),
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
	}
	datasetID := viper.GetString("dataset")
	assemblyID := viper.GetString("assembly")
	return vcf.IterateOver(r, datasetID, assemblyID, func(v *variant.Variant) error {
		return c.InsertVariant(v)
	})
}
