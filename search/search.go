package search

import (
	"github.com/labbcb/brave/variant"
)

// Response is the response object
type Response struct {
	Draw            int                `json:"draw"`            // the draw counter that this object is a response to
	RecordsTotal    int                `json:"recordsTotal"`    // total records, before filtering
	RecordsFiltered int                `json:"recordsFiltered"` // total records, after filtering
	Error           string             `json:"error"`           // list of variants that matched one (or more) search query
	Variants        []*variant.Variant `json:"data"`            // if an error occurs during the running of the server-side processing
}

// Input it the request object
type Input struct {
	Draw    int      `json:"draw"`    // draw counter
	Start   int64    `json:"start"`   // paging first record indicator
	Length  int64    `json:"length"`  // number of records that the table can display in the current draw
	Queries []*Query `json:"queries"` // list of queries
}
