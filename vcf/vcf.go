package vcf

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/brentp/vcfgo"
	"github.com/labbcb/brave/variant"
)

const (
	// AF is the allele frequency for each ALT allele in the same order as listed (estimated from primary data, not called genotypes)
	AF = "AF"
	// NS is the number of samples with data
	NS = "NS"
	// ANN is an INFO column related to variant annotation
	ANN = "ANN"
	// CLNSIG is an INFO column related to variant annotation
	CLNSIG = "CLNSIG"
	// DP is per-sample read depth
	DP = "DP"
	// GQ is per-sample conditional genotype quality
	GQ = "GQ"
	// GeneSymbol is an ANN field for gene symbol, one per alternate bases
	GeneSymbol = 3
	// Type is an ANN field for variant type, one per alternate bases
	Type = 5
	// HGVS is an ANN field for HGVS nomenclature, one per alternate bases
	HGVS = 9
)

type VCFSummary struct {
	TotalVariants  uint
	PassedVariants uint
}

// IterateOver reads a VCF file (from io.Reader) and yeld varint to caller's function.
// It returns VCFSummary with total of variants read and total of variants that passed all filters (FILTER = PASS or .).
// If err is not nil, VCFSummary will have the total variants so far.
func IterateOver(r io.Reader, doFilter bool, datasetID, assemblyID string, doSomething func(v *variant.Variant) error) (VCFSummary, error) {
	vcfReader, err := vcfgo.NewReader(r, false)
	if err != nil {
		return VCFSummary{}, err
	}
	totalSamples := len(vcfReader.Header.SampleNames)

	buildVariant := func(v *vcfgo.Variant) *variant.Variant {
		return &variant.Variant{
			DatasetID:       datasetID,
			TotalSamples:    int32(totalSamples),
			AssemblyID:      assemblyID,
			SnpIds:          getSnpIds(v),
			ReferenceName:   strings.TrimPrefix(v.Chromosome, "chr"),
			Start:           int32(v.Pos),
			ReferenceBases:  v.Reference,
			AlternateBases:  v.Alternate,
			GeneSymbol:      GetAnnotationColumn(v, GeneSymbol),
			AlleleFrequency: GetAttributeAsFloatSlice(v, AF, nil),
			SampleCount:     GetAttributeAsInt(v, NS, 0),
			Coverage:        CalculateDistribution(GetSamplesDP(v)),
			GenotypeQuality: CalculateDistribution(GetSamplesGQ(v)),
			CLNSIG:          GetAttributeAsString(v, CLNSIG, ""),
			HGVS:            GetAnnotationColumn(v, HGVS),
			Type:            GetAnnotationColumn(v, Type),
		}
	}

	var totalVariants, passedVariants uint

	for {
		v := vcfReader.Read()
		if v == nil {
			break
		}

		totalVariants += 1

		if doFilter && v.Filter != "PASS" && v.Filter != "." {
			continue
		}

		if err := doSomething(buildVariant(v)); err != nil {
			return VCFSummary{totalVariants, passedVariants}, err
		}

		passedVariants += 1
	}
	return VCFSummary{totalVariants, passedVariants}, vcfReader.Error()
}

func getSnpIds(v *vcfgo.Variant) []string {
	if v.Id_ == "." {
		return nil
	}
	return strings.Split(v.Id_, ";")
}

// GetAttributeAsString gets an INFO field and return as string
func GetAttributeAsString(v *vcfgo.Variant, key string, defaultValue string) string {
	i, _ := v.Info_.Get(key)
	switch s := i.(type) {
	case string:
		return s
	case []string:
		return strings.Join(s, ",")
	default:
		return defaultValue
	}
}

// GetAnnotationColumn gets an ANN column by its index
func GetAnnotationColumn(v *vcfgo.Variant, index int) []string {
	i, _ := v.Info_.Get(ANN)
	if i == nil {
		return nil
	}

	switch ann := i.(type) {
	case string:
		var columns []string
		for _, a := range strings.Split(ann, ",") {
			columns = append(columns, strings.Split(a, "|")[index])
		}
		return columns
	case []string:
		var columns []string
		for _, a := range ann {
			columns = append(columns, strings.Split(a, "|")[index])
		}
		return columns
	default:
		panic(fmt.Sprintf("invalid type %T", ann))
	}
}

// GetAttributeAsFloatSlice gets AF values
func GetAttributeAsFloatSlice(v *vcfgo.Variant, key string, defaultValue []float32) []float32 {
	i, err := v.Info_.Get(key)
	if err != nil {
		return defaultValue
	}
	return i.([]float32)
}

// GetAttributeAsInt gets INFO key and return as int
func GetAttributeAsInt(v *vcfgo.Variant, key string, defaultValue int) int {
	i, err := v.Info_.Get(key)
	if err != nil {
		return defaultValue
	}
	return i.(int)
}

// GetSamplesDP get per-sample DP
func GetSamplesDP(v *vcfgo.Variant) []int {
	var dps []int
	for _, s := range v.Samples {
		dps = append(dps, s.DP)
	}
	return dps
}

// GetSamplesGQ per-sample GQ
func GetSamplesGQ(v *vcfgo.Variant) []int {
	var gqs []int
	for _, s := range v.Samples {
		gqs = append(gqs, s.GQ)
	}
	return gqs
}

// CalculateDistribution calculates distribution of a given slice of int
func CalculateDistribution(xs []int) *variant.Distribution {
	length := len(xs)

	// if there is no sample
	if length == 0 {
		return nil
	}

	// if there is only one sample
	if length == 1 {
		x := xs[0]
		return &variant.Distribution{
			Min:    float64(x),
			Q25:    float64(x),
			Median: float64(x),
			Q75:    float64(x),
			Max:    float64(x),
			Mean:   float64(x),
		}
	}

	// sort values
	sort.Ints(xs)

	// calculate median
	half := length / 2
	var median float64
	if length%2 == 0 {
		median = float64(xs[half]+xs[half-1]) / 2.0
	} else {
		median = float64(xs[half])
	}

	// calculate sum
	var sum int
	for _, x := range xs {
		sum += x
	}

	// calculate q25, q75 and mean
	return &variant.Distribution{
		Min:    float64(xs[0]),
		Q25:    float64(xs[int(0.25*(float64(length)+1))-1]),
		Median: median,
		Q75:    float64(xs[int(0.75*(float64(length)+1))-1]),
		Max:    float64(xs[length-1]),
		Mean:   float64(sum) / float64(length),
	}
}
