package stats

import (
	"log"
)

// Aggregator is a type that provides a method to aggregate values gathered
// by metric when it has queried the DB
type Aggregator interface {
	Aggregate(output []float64, values []float64) []float64
}

// NewAggregator returns an aggregator that matches string a
func NewAggregator(a string) Aggregator {
	switch a {
	case "sum":
		return &Sum{}
	case "accumulate":
		return &Accumulate{}
	}
	log.Fatal("unknown aggregator " + a)
	return nil
}

// Sum is an aggregator that sums the values provided and appends it to output
type Sum struct{}

// Aggregate aggregates values provided into output
func (a *Sum) Aggregate(output []float64, values []float64) []float64 {
	var sum float64
	for i := range values {
		sum += values[i]
	}
	return append(output, sum)
}

// Accumulate is an aggregator that sums the values provided and the last value
// in the output, which makes for an accumulated value over time
type Accumulate struct{}

// Aggregate aggregates values provided into output
func (a *Accumulate) Aggregate(output []float64, values []float64) []float64 {
	var sum float64
	for i := range values {
		sum += values[i]
	}
	if len(output) > 0 {
		sum += output[len(output)-1]
	}
	return append(output, sum)
}
