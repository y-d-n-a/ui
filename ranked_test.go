
package ui

import (
	"strings"
	"testing"
)

// ranked is extracted here as a pure function for testing.
func rankedOptions(options []Option, query string) []Option {
	if query == "" {
		out := make([]Option, len(options))
		copy(out, options)
		return out
	}
	q := strings.ToLower(query)
	type scored struct {
		opt   Option
		score int
		idx   int
	}
	var items []scored
	for i, o := range options {
		lower := strings.ToLower(o.Label)
		var score int
		if strings.HasPrefix(lower, q) {
			score = 0
		} else if strings.Contains(lower, q) {
			score = 1
		} else {
			score = 2
		}
		items = append(items, scored{o, score, i})
	}
	// stable sort by score then original index
	for i := 1; i < len(items); i++ {
		for j := i; j > 0; j-- {
			a, b := items[j-1], items[j]
			if a.score > b.score || (a.score == b.score && a.idx > b.idx) {
				items[j-1], items[j] = items[j], items[j-1]
			}
		}
	}
	out := make([]Option, len(items))
	for i, it := range items {
		out[i] = it.opt
	}
	return out
}

var testOptions = []Option{
	{Label: "Apple", Value: "apple"},
	{Label: "Banana", Value: "banana"},
	{Label: "Apricot", Value: "apricot"},
	{Label: "Cherry", Value: "cherry"},
	{Label: "Pineapple", Value: "pineapple"},
}

func TestRanked_EmptyQuery(t *testing.T) {
	result := rankedOptions(testOptions, "")
	if len(result) != len(testOptions) {
		t.Fatalf("expected %d results, got %d", len(testOptions), len(result))
	}
	for i, o := range result {
		if o.Value != testOptions[i].Value {
			t.Errorf("position %d: expected %q, got %q", i, testOptions[i].Value, o.Value)
		}
	}
}

func TestRanked_PrefixFirst(t *testing.T) {
	result := rankedOptions(testOptions, "ap")
	// Apple and Apricot prefix-match "ap", Pineapple contains "ap"
	if len(result) == 0 {
		t.Fatal("expected results")
	}
	if result[0].Value != "apple" {
		t.Errorf("expected apple first, got %q", result[0].Value)
	}
	if result[1].Value != "apricot" {
		t.Errorf("expected apricot second, got %q", result[1].Value)
	}
}

func TestRanked_ContainsAfterPrefix(t *testing.T) {
	result := rankedOptions(testOptions, "apple")
	// "apple" prefix matches Apple; Pineapple contains "apple"
	if result[0].Value != "apple" {
		t.Errorf("expected apple first, got %q", result[0].Value)
	}
	if result[1].Value != "pineapple" {
		t.Errorf("expected pineapple second, got %q", result[1].Value)
	}
}

func TestRanked_NoMatch(t *testing.T) {
	result := rankedOptions(testOptions, "zzz")
	// All score 2 (no match), order preserved
	for i, o := range result {
		if o.Value != testOptions[i].Value {
			t.Errorf("position %d: expected %q, got %q", i, testOptions[i].Value, o.Value)
		}
	}
}

func TestRanked_CaseInsensitive(t *testing.T) {
	result := rankedOptions(testOptions, "APPLE")
	if result[0].Value != "apple" {
		t.Errorf("expected apple first for uppercase query, got %q", result[0].Value)
	}
}
