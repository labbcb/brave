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

## Environment variables

Server specific.

- `BRAVE_DATABASE` URL to Mongo database. Default is `mongodb://localhost:27017`.
- `BRAVE_ADDRESS` address to bind server. Default is `:8080`.
- `BRAVE_USERNAME` administrator user name. Default is `admin`.
- `BRAVE_PASSWORD` administrator password. Default is empty.


## Deploy server with Docker

Create network and volume.

```bash
docker volume create brave-data
docker network create brave

docker container run \
    --network brave \
    --rm --detach \
    --name brave-db \
    --volume brave-data:/data/db \
    mongo:4

docker container run \
    --rm --detach \
    --name brave \
    --network brave \
    --publish 8080:8080 \
    --env BRAVE_DATABASE=mongodb://brave2-db:27017 \
    --env BRAVE_PASSWORD=secret \
    bipmed/brave server
```

## Import variants

BraVE accepts VCF files (v4.2) as input and submit variants to server instance. No genotype (FORMAT column) data is sent to server. FORMAT/DP and FORMAT/GQ are used to calculate distribution (min, q25, median, q75, max and average) of every variant. By default only variant that passed all filters are imported to database (FILTER = PASS or .). Use `--dont-filter` option to import all variants, regardless of FILTER column.

```bash
brave import \
    [--dont-filter] \
    [--dry-run] \
    [--host http://localhost:8080] \
    [--username admin] \
    --password secret \
    --assembly hg38 \
    --dataset bipmed \
    bipmed.hg38.vcf.gz
```