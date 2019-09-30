package search

import (
	"regexp"
	"strconv"
)

var (
	// GenomicRange is a regex that matches a genomic range, 1:1000-2000
	GenomicRange = regexp.MustCompile(`^\s*([1-9]|1[0-9]|2[0-2]|[XY])\s*:\s*(\d+)\s*-\s*(\d+)\s*$`)
	// GenomicPosition is a regex that matches a genomic position, 1:1000
	GenomicPosition = regexp.MustCompile(`^\s*([1-9]|1[0-9]|2[0-2]|[XY])\s*:\s*(\d+)\s*$`)
	// SnpID is a regex that matches dbSNP ID, rs35735053
	SnpID = regexp.MustCompile(`^\s*(rs\d+)\s*$`)
	// GeneSymbol is a regex that matches gene name, SCN1A
	GeneSymbol = regexp.MustCompile(`^\s*([A-Za-z0-9\-]+)\s*$`)
)

// Query contains optional parameters for filtering variants.
// All fields may be omitted meaning that matches with all variants present in the database
type Query struct {
	SnpID         string `json:"snpId"`         // external variant id, normally from dbSNP database (rs35735053)
	AssemblyID    string `json:"assemblyId"`    // reference genome version (GRCh38)
	DatasetID     string `json:"datasetId"`     // call set id (bipmed-wes-phase2)
	ReferenceName string `json:"referenceName"` // chromosome name (chr1, 1)
	Start         int32  `json:"start"`         // start position (7737651)
	End           int32  `json:"end"`           // end position (70000)
	GeneSymbol    string `json:"geneSymbol"`    // gene symbol (SCN1A)
}

// Parse parses text to a query
func Parse(text string) *Query {
	xs := GenomicRange.FindStringSubmatch(text)
	if xs != nil {
		return &Query{ReferenceName: xs[1], Start: mustBeInt32(xs[2]), End: mustBeInt32(xs[3])}
	}
	xs = GenomicPosition.FindStringSubmatch(text)
	if xs != nil {
		return &Query{ReferenceName: xs[1], Start: mustBeInt32(xs[2])}
	}
	xs = SnpID.FindStringSubmatch(text)
	if xs != nil {
		return &Query{SnpID: xs[1]}
	}
	xs = GeneSymbol.FindStringSubmatch(text)
	if xs != nil {
		return &Query{GeneSymbol: xs[1]}
	}
	return new(Query)
}

func mustBeInt32(s string) int32 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(`query: ParseUint("` + s + `"): ` + err.Error())
	}
	return int32(i)
}
