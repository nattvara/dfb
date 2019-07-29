package stats

import (
	"log"
)

// Aggregators is a map of availible aggregators
var Aggregators = map[string]Aggregator{
	"sum":        &Sum{},
	"accumulate": &Accumulate{},
}

// Aggregator is a type that provides a method to aggregate values gathered
// by metric when it has queried the DB
type Aggregator interface {
	Aggregate(output []float64, values []float64) []float64
}

// NewAggregator returns an aggregator that matches string name
func NewAggregator(name string) Aggregator {
	if _, ok := Aggregators[name]; !ok {
		log.Fatal("unknown aggregator " + name)
	}
	a := Aggregators[name]
	return a
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
