package stats

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-memdb"
)

const (
	// TimeUnitDays is an identifier that means a metrics values should be grouped by day
	TimeUnitDays = "days"

	// TimeUnitMonths is an identifier that means a metrics values should be grouped by month
	TimeUnitMonths = "months"

	// TimeUnitYears is an identifier that means a metrics values should be grouped by year
	TimeUnitYears = "years"
)

// TimeUnits contains the supported time units
var TimeUnits = []string{
	TimeUnitDays,
	TimeUnitMonths,
	TimeUnitYears,
}

// getDateLayoutForTimeUnit returns the layout used for stringifying a time.Time for given timeUnit
func getDateLayoutForTimeUnit(timeUnit string) string {
	var layout string
	switch strings.ToLower(timeUnit) {
	case TimeUnitDays:
		layout = "2006-01-02"
	case TimeUnitMonths:
		layout = "Jan 2006"
	case TimeUnitYears:
		layout = "2006"
	default:
		log.Fatal("unsupported timeUnit: " + timeUnit)
	}
	return layout
}

// DateIterator is a type used for iterating over a range of dates and querying
// a DB for records for a given date
type DateIterator struct {
	TimeUnit      string
	TimeLength    int
	CurrentOffset int

	db          *DB
	currentDate *Date
}

// Date is a type containing a time.Time Value and boolean Valid fields,
// Valid is true the Value is inside the allowed range of dates
type Date struct {
	Value time.Time
	Valid bool
}

// NewDateIterator returns a new date iterator
func NewDateIterator(db *DB, timeUnit string, timeLength int) DateIterator {
	it := DateIterator{
		TimeUnit:   timeUnit,
		TimeLength: timeLength,
		db:         db,
	}
	it.CurrentOffset = -1
	it.currentDate = &Date{Value: time.Now(), Valid: true}
	it.decrementDate(timeLength + 1)
	return it
}

// Next returns the next date after currentDate of DateIterator it. Also updates
// Valid field accordingly
func (it *DateIterator) Next() *Date {
	it.CurrentOffset++
	if it.CurrentOffset > it.TimeLength {
		it.currentDate.Valid = false
	}
	it.incrementDate(1)
	return it.currentDate
}

// GetRecordsForDate queries the DB of DateIterator it for values in
// provided table matching date and metric metadata and returns matching records
func (it *DateIterator) GetRecordsForDate(table string, metric Metric, date time.Time) memdb.ResultIterator {
	layout := getDateLayoutForTimeUnit(it.TimeUnit)

	txn := it.db.memdb.Txn(false)
	defer txn.Abort()

	var args []interface{}
	args = append(args, metric.GetMetadata("repo"))
	args = append(args, metric.GetMetadata("group"))
	if metric.SupportsDomains() {
		args = append(args, metric.GetMetadata("domain"))
	}
	args = append(args, date.Format(layout))

	records, err := txn.Get(
		table,
		it.db.GetIndexFromTimeUnit(it.TimeUnit, metric.SupportsDomains()),
		args...,
	)
	if err != nil {
		log.Fatal("failed to fetch data from db for metric", err)
	}

	return records
}

// incrementDate increments currentDate of DateIterator it by given amount
func (it *DateIterator) incrementDate(amount int) {
	switch strings.ToLower(it.TimeUnit) {
	case TimeUnitDays:
		it.currentDate.Value = it.currentDate.Value.AddDate(0, 0, amount)
	case TimeUnitMonths:
		it.currentDate.Value = it.currentDate.Value.AddDate(0, amount, 0)
	case TimeUnitYears:
		it.currentDate.Value = it.currentDate.Value.AddDate(amount, 0, 0)
	default:
		log.Fatal("unsupported time unit: " + it.TimeUnit)
	}
}

// decrementDate decrements currentDate of DateIterator it by given amount
func (it *DateIterator) decrementDate(amount int) {
	amount *= -1
	switch strings.ToLower(it.TimeUnit) {
	case TimeUnitDays:
		it.currentDate.Value = it.currentDate.Value.AddDate(0, 0, amount)
	case TimeUnitMonths:
		it.currentDate.Value = it.currentDate.Value.AddDate(0, amount, 0)
	case TimeUnitYears:
		it.currentDate.Value = it.currentDate.Value.AddDate(amount, 0, 0)
	default:
		log.Fatal("unsupported time unit: " + it.TimeUnit)
	}
}
