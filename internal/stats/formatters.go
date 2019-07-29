package stats

import (
	"fmt"
	"math"
)

// Formatter is a type that implements a Format method for a float64 to a string
type Formatter interface {
	Format(value float64) string
}

// BytesFormatter formats a float64 of bytes to a  a nicely formatted
// string eg. 289.0 MiB 3.1 GiB
type BytesFormatter struct{}

// Format formats provided float64 value to a string
// source: https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func (f *BytesFormatter) Format(value float64) string {
	bytes := value
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%v B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf(
		"%.1f %ciB",
		float64(bytes)/float64(div),
		"KMGTPE"[exp],
	)
}

// AmountFormatter formats an amount of "things" to a shorter representation as
// a string, eg 1000 -> 1k, 1000000 -> 1m, etc.
type AmountFormatter struct{}

// Format formats provided float64 value to a string
func (f *AmountFormatter) Format(value float64) string {
	var exp string
	if value >= 1000000000 {
		value = math.Round(value*10) / 10000000000
		exp = "b"
	} else if value >= 1000000 {
		value = math.Round(value*10) / 10000000
		exp = "m"
	} else if value >= 1000 {
		value = math.Round(value*10) / 10000
		exp = "k"
	}
	return fmt.Sprintf("%.1f%s", value, exp)
}
