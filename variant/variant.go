package variant

import "fmt"

// Variant is a genomic variant that was annotated, sample data removed and calculated distribution.
// Variants types supported by VCF are: Integer (32-bit, signed), Float (32-bit, IEEE-754).
type Variant struct {
	ID              string        `json:"id" bson:"_id"`                                    // variant id
	DatasetID       string        `json:"datasetId" bson:"datasetId"`                       // dataset ID
	TotalSamples    int32         `json:"totalSamples" bson:"totalSamples"`                 // total samples in dataset
	AssemblyID      string        `json:"assemblyId" bson:"assemblyId"`                     // reference genome version (b37, hg38)
	SnpIds          []string      `json:"snpIds,omitempty" bson:"snpIds"`                   // ids (ID)
	ReferenceName   string        `json:"referenceName" bson:"referenceName"`               // contig name (CHROM)
	Start           int32         `json:"start"`                                            // 0-based position (POS)
	ReferenceBases  string        `json:"referenceBases,omitempty" bson:"referenceBases"`   // reference bases (REF)
	AlternateBases  []string      `json:"alternateBases,omitempty" bson:"alternateBases"`   // list of alternate bases (ALT)
	GeneSymbol      []string      `json:"geneSymbol,omitempty" bson:"geneSymbol"`           // gene symbol, one per ALT
	AlleleFrequency []float32     `json:"alleleFrequency,omitempty" bson:"alleleFrequency"` // allele frequency (AF), one per ALT
	SampleCount     int           `json:"sampleCount" bson:"sampleCount"`                   // total samples that have this variant (NS)
	Coverage        *Distribution `json:"coverage,omitempty"`                               // distribution of coverage (DP)
	GenotypeQuality *Distribution `json:"genotypeQuality,omitempty" bson:"genotypeQuality"` //distribution of genotype quality (GQ)
	CLNSIG          string        `json:"clnsig"`                                           // clinical significance
	HGVS            []string      `json:"hgvs,omitempty"`                                   // HGVS nomenclature
	Type            []string      `json:"type,omitempty"`                                   // variant type
}

// Distribution represents distribution of a list of values.
type Distribution struct {
	Min    float64 `json:"min"`    // minimum value
	Q25    float64 `json:"q25"`    // 25% percentile
	Median float64 `json:"median"` // median
	Q75    float64 `json:"q75"`    // 75% percentile
	Max    float64 `json:"max"`    // maximum value
	Mean   float64 `json:"mean"`   // average
}

func (v *Variant) String() string {
	return fmt.Sprintf("%s-%s %s:%d (%d/%d) %v %s > %v (AF=%v DP={%s} GQ={%s} GENES=%v CLNSIG=%v HGVS=%s TYPE=%v)",
		v.DatasetID,
		v.AssemblyID,
		v.ReferenceName,
		v.Start,
		v.SampleCount,
		v.TotalSamples,
		v.SnpIds,
		v.ReferenceBases,
		v.AlternateBases,
		v.AlleleFrequency,
		v.Coverage.String(),
		v.GenotypeQuality.String(),
		v.GeneSymbol,
		v.CLNSIG,
		v.HGVS,
		v.Type)
}

func (d *Distribution) String() string {
	if d == nil {
		return ""
	}
	return fmt.Sprintf("min=%.2f q25=%.2f median=%.2f q75=%.2f max=%.2f mean=%.2f",
		d.Min,
		d.Q25,
		d.Median,
		d.Q75,
		d.Max,
		d.Mean)
}
