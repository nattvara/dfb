package stats

import "errors"

// Aggregators is a map of availible aggregators
var Aggregators = map[string]Aggregator{
	"sum":        &Sum{},
	"accumulate": &Accumulate{},
	"average":    &Average{},
}

// Aggregator is a type that provides a method to aggregate values gathered
// by metric when it has queried the DB
type Aggregator interface {
	Aggregate(output []float64, values []float64) []float64
}

// NewAggregator returns an aggregator that matches string name
func NewAggregator(name string) (Aggregator, error) {
	if _, ok := Aggregators[name]; !ok {
		return nil, errors.New("unknown aggregator " + name)
	}
	a := Aggregators[name]
	return a, nil
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

// Average is an aggregator that averages the values provided and appends it to output
type Average struct{}

// Aggregate aggregates values provided into output
func (a *Average) Aggregate(output []float64, values []float64) []float64 {
	var sum float64
	var avg float64
	number := float64(len(values))

	for i := range values {
		sum += values[i]
	}
	if number > 0 {
		avg = sum / float64(len(values))
	}
	return append(output, avg)
}
