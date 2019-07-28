package stats

import (
	"fmt"
	"log"
	"reflect"
	"time"
)

const (
	// AllDomains is a special string to mean all domains (".." cannot be a valid domain name)
	AllDomains = ".."
)

var metrics = map[string]Metric{
	"snapshots-data-added":            &SnapshotsDataAdded{},
	"snapshots-files-new-and-changed": &SnapshotsFilesNewAndChanged{},
	"snapshots-files-processed":       &SnapshotsFilesProcessed{},
}

// NewMetric returns a new metric
func NewMetric(name string, repo string, group string, domain string, timeUnit string, aggregator string) Metric {
	if _, ok := metrics[name]; !ok {
		log.Fatal("unknown metric " + name)
	}
	m := metrics[name]
	m.Init(timeUnit)
	m.SetTitle(name, repo, group, domain, aggregator)
	m.SetMetadata(name, repo, group, domain, aggregator)
	return m
}

// Metric is a type that provides a FetchDataFromDB method and methods for
// retrieving values, labels, and titles
type Metric interface {
	SetTitle(name string, repo string, group string, domain string, aggregator string)
	GetTitle() string
	SetMetadata(name string, repo string, group string, domain string, aggregator string)
	GetValues(a Aggregator) []float64
	GetLabels() []time.Time
	GetDateLayout() string
	GetFormatter() Formatter

	Init(timeUnit string)
	FetchDataFromDB(db *DB, timeUnit string, timeLength int)
}

// metricData is type that provides for setter and getters for Metrics
type metricData struct {
	Title      string
	Name       string
	Meta       map[string]string
	Data       [][]float64
	Dates      []time.Time
	DateLayout string
	Formatter  Formatter
}

// SetTitle sets the title of a metricData m from given input data
func (m *metricData) SetTitle(name string, repo string, group string, domain string, aggregator string) {
	var d string
	if domain == AllDomains {
		d = "all domains"
	} else {
		d = domain
	}
	m.Title = fmt.Sprintf(
		"%s %s to %s in group %s of repo %s",
		aggregator,
		m.Name,
		d,
		group,
		repo,
	)
}

// GetTitle returns the title of metricData m
func (m *metricData) GetTitle() string {
	return m.Title
}

// SetMetadata sets the metadata of metricData m
func (m *metricData) SetMetadata(name string, repo string, group string, domain string, aggregator string) {
	m.Meta = map[string]string{"repo": repo, "group": group, "domain": domain, "aggregator": aggregator}
}

// AddDate adds a date that metricData m should have values for
func (m *metricData) AddDate(date time.Time) {
	m.Dates = append(m.Dates, date)
	m.Data = append(m.Data, []float64{})
}

// AppendValues appends values from given field of provided object to given date for metricData m
func (m *metricData) AppendValues(obj interface{}, field string, date int) {
	v := reflect.ValueOf(obj)
	v = reflect.Indirect(v)
	fv := v.FieldByName(field)

	var value float64
	switch fv.Kind().String() {
	case "float":
		value = fv.Float()
	case "int":
		value = float64(fv.Int())
	default:
		panic("unkown field value " + fv.Kind().String())
	}

	m.Data[date] = append(m.Data[date], value)
}

// AppendMultipleValues appends multiple fields to metricData m
func (m *metricData) AppendMultipleValues(obj interface{}, fields []string, date int) {
	for _, field := range fields {
		m.AppendValues(obj, field, date)
	}
}

// GetValues returns aggregated values from metricData m
func (m *metricData) GetValues(a Aggregator) []float64 {
	var res []float64
	for _, values := range m.Data {
		res = a.Aggregate(res, values)
	}
	return res
}

// GetLabels returns labels (time.Time) for values of metricData m
func (m *metricData) GetLabels() []time.Time {
	return m.Dates
}

// GetDateLayout returns the layout used for dates for metricData m
func (m *metricData) GetDateLayout() string {
	return m.DateLayout
}

// GetFormatter returns the formatter metricData m uses to format its values
func (m *metricData) GetFormatter() Formatter {
	return m.Formatter
}

// SnapshotsDataAdded is a metric of the data added by snapshots over time
type SnapshotsDataAdded struct {
	metricData
}

// Init initializes the metric
func (m *SnapshotsDataAdded) Init(timeUnit string) {
	m.DateLayout = getDateLayoutForTimeUnit(timeUnit)
	m.Name = "data added"
	m.Formatter = &BytesFormatter{}
}

// FetchDataFromDB from db fetches appropriate data from DB and appends values
// for given timeUnit and timeLength
func (m *SnapshotsDataAdded) FetchDataFromDB(db *DB, timeUnit string, timeLength int) {
	iterator := NewDateIterator(db, timeUnit, timeLength)
	for date := iterator.Next(); date.Valid; date = iterator.Next() {
		records := iterator.QueryDB("snapshot", m.Meta, date.Value)
		m.AddDate(date.Value)
		for obj := records.Next(); obj != nil; obj = records.Next() {
			m.AppendValues(obj, "DataAdded", iterator.CurrentOffset)
		}
	}
}

// SnapshotsFilesNewAndChanged is a metric of the new and changed files for the snapshots over time
type SnapshotsFilesNewAndChanged struct {
	metricData
}

// Init initializes the metric
func (m *SnapshotsFilesNewAndChanged) Init(timeUnit string) {
	m.DateLayout = getDateLayoutForTimeUnit(timeUnit)
	m.Name = "new files"
	m.Formatter = &AmountFormatter{}
}

// FetchDataFromDB from db fetches appropriate data from DB and appends values
// for given timeUnit and timeLength
func (m *SnapshotsFilesNewAndChanged) FetchDataFromDB(db *DB, timeUnit string, timeLength int) {
	iterator := NewDateIterator(db, timeUnit, timeLength)
	for date := iterator.Next(); date.Valid; date = iterator.Next() {
		records := iterator.QueryDB("snapshot", m.Meta, date.Value)
		m.AddDate(date.Value)
		for obj := records.Next(); obj != nil; obj = records.Next() {
			m.AppendMultipleValues(obj, []string{"FilesNew", "FilesChanged"}, iterator.CurrentOffset)
		}
	}
}

// SnapshotsFilesProcessed is a metric of the files processed for the snapshots over time
type SnapshotsFilesProcessed struct {
	metricData
}

// Init initializes the metric
func (m *SnapshotsFilesProcessed) Init(timeUnit string) {
	m.DateLayout = getDateLayoutForTimeUnit(timeUnit)
	m.Name = "new files"
	m.Formatter = &AmountFormatter{}
}

// FetchDataFromDB from db fetches appropriate data from DB and appends values
// for given timeUnit and timeLength
func (m *SnapshotsFilesProcessed) FetchDataFromDB(db *DB, timeUnit string, timeLength int) {
	iterator := NewDateIterator(db, timeUnit, timeLength)
	for date := iterator.Next(); date.Valid; date = iterator.Next() {
		records := iterator.QueryDB("snapshot", m.Meta, date.Value)
		m.AddDate(date.Value)
		for obj := records.Next(); obj != nil; obj = records.Next() {
			m.AppendValues(obj, "TotalFilesProcessed", iterator.CurrentOffset)
		}
	}
}
