# BIPMed Variant Explorer

This is the reference implementation server of the [Brave API](api/swagger.yaml).
It accepts `Query` as POST request and returns a list of `Variant`.

Query

- `variantId` - Unique identifier (rs6054257)
- `assemblyId` - Version of reference genomes (GCRh38, `null` for all reference genomes)
- `datasetId` - Dataset id (bipmedExome, `null` for all datasets)
- `referenceName` - Chromosome (20, must be present when `geneSymbol` is `null`)
- `start` - Exact position when `end` is `null`, otherwise is start (include) of range (1000, must be present)
- `end` - End (include) of range (must be present when `start` is not `null`)
- `geneSymbol` - Gene Symbol (SCN1A, it can be combined with `referenceGenome`, `start` and `end`)

Variant

- `variantIds` - List of unique identifiers where available, can be `null` (ID)
- `assemblyId` - Version of reference genome, required
- `datasetId` - Dataset id, required
- `referenceName` - An identifier from the reference genome, required (CHROM)
- `start` - The reference position, with the 1st base having position 1, required (POS)
- `referenceBases` - Reference bases, required (REF)
- `alternateBases` - List of alternate non-reference alleles, required (ALT)
- `geneSymbol` - Gene symbol, can be `null`
- `alleleFrequency` - Allele Frequency, required (AF)
- `sampleCount` - Number of Samples With Data, required (NS)
