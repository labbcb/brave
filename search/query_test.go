package search

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	ts := map[string]*Query{
		"SCN1A":         {GeneSymbol: "SCN1A"},
		"1:65000-70000": {ReferenceName: "1", Start: 65000, End: 70000},
		"1:7737651":     {ReferenceName: "1", Start: 7737651},
		"rs35735053":    {SnpID: "rs35735053"},
		"":              {},
	}

	for text, want := range ts {
		got := Parse(text)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %v, got %v", want, got)
		}
	}
}
