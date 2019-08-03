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

// Metrics is a map of availible metrics
var Metrics = map[string]Metric{
	"backup-time":                  &BackupTime{},
	"repo-disk-space":              &RepoDiskSpace{},
	"domain-data-added":            &DomainDataAdded{},
	"domain-disk-space":            &DomainDiskSpace{},
	"domain-disk-space-on-restore": &DomainDiskSpaceOnRestore{},
	"domain-files-new-and-changed": &DomainFilesNewAndChanged{},
	"domain-files-processed":       &DomainFilesProcessed{},
}

// NewMetric returns a new metric
func NewMetric(name string, repo string, group string, domain string, timeUnit string, aggregator string) Metric {
	if _, ok := Metrics[name]; !ok {
		log.Fatal("unknown metric " + name)
	}
	m := Metrics[name]
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
	SupportsDomains() bool
	SetMetadata(name string, repo string, group string, domain string, aggregator string)
	GetMetadata(property string) string
	GetValues(a Aggregator) []float64
	GetLabels() []time.Time
	GetDateLayout() string
	GetFormatter() Formatter

	Init(timeUnit string)
	FetchDataFromDB(db *DB, timeUnit string, timeLength int)
}

// metricData is type that provides for setter and getters for Metrics
type metricData struct {
	Title           string
	Name            string
	supportsDomains bool
	Meta            map[string]string
	Data            [][]float64
	Dates           []time.Time
	DateLayout      string
	Formatter       Formatter
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

// SupportsDomains returns whether metricData m supports specifying domain
func (m *metricData) SupportsDomains() bool {
	return m.supportsDomains
}

// SetMetadata sets the metadata of metricData m
func (m *metricData) SetMetadata(name string, repo string, group string, domain string, aggregator string) {
	m.Meta = map[string]string{"repo": repo, "group": group, "domain": domain, "aggregator": aggregator}
}

// GetMetadata returns value for requested metadata property from metricData m
func (m *metricData) GetMetadata(property string) string {
	if _, ok := m.Meta[property]; !ok {
		log.Fatal("unknown metadata property " + property)
	}
	return m.Meta[property]
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
	case "float64":
		value = fv.Float()
	case "int":
		value = float64(fv.Int())
	case "int64":
		value = float64(fv.Int())
	default:
		panic("cannot append value, unkown type")
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

// DomainDataAdded is a metric of the data added by snapshots of a domain over time
type DomainDataAdded struct {
	metricData
}

// Init initializes the metric
func (m *DomainDataAdded) Init(timeUnit string) {
	m.supportsDomains = true
	m.DateLayout = getDateLayoutForTimeUnit(timeUnit)
	m.Name = "data added"
	m.Formatter = &BytesFormatter{}
}

// FetchDataFromDB fetches appropriate data from DB and appends values
// for given timeUnit and timeLength
func (m *DomainDataAdded) FetchDataFromDB(db *DB, timeUnit string, timeLength int) {
	iterator := NewDateIterator(db, timeUnit, timeLength)
	for date := iterator.Next(); date.Valid; date = iterator.Next() {
		records := iterator.GetRecordsForDate("snapshot", m, date.Value)
		m.AddDate(date.Value)
		for obj := records.Next(); obj != nil; obj = records.Next() {
			m.AppendValues(obj, "DataAdded", iterator.CurrentOffset)
		}
	}
}

// DomainFilesNewAndChanged is a metric of the new and changed files for the snapshots
// of a domain over time
type DomainFilesNewAndChanged struct {
	metricData
}

// Init initializes the metric
func (m *DomainFilesNewAndChanged) Init(timeUnit string) {
	m.supportsDomains = true
	m.DateLayout = getDateLayoutForTimeUnit(timeUnit)
	m.Name = "new and changed files"
	m.Formatter = &AmountFormatter{}
}

// FetchDataFromDB fetches appropriate data from DB and appends values
// for given timeUnit and timeLength
func (m *DomainFilesNewAndChanged) FetchDataFromDB(db *DB, timeUnit string, timeLength int) {
	iterator := NewDateIterator(db, timeUnit, timeLength)
	for date := iterator.Next(); date.Valid; date = iterator.Next() {
		records := iterator.GetRecordsForDate("snapshot", m, date.Value)
		m.AddDate(date.Value)
		for obj := records.Next(); obj != nil; obj = records.Next() {
			m.AppendMultipleValues(obj, []string{"FilesNew", "FilesChanged"}, iterator.CurrentOffset)
		}
	}
}

// DomainFilesProcessed is a metric of the files processed while taking a
// snapshot of a domain
type DomainFilesProcessed struct {
	metricData
}

// Init initializes the metric
func (m *DomainFilesProcessed) Init(timeUnit string) {
	m.supportsDomains = true
	m.DateLayout = getDateLayoutForTimeUnit(timeUnit)
	m.Name = "files processed"
	m.Formatter = &AmountFormatter{}
}

// FetchDataFromDB fetches appropriate data from DB and appends values
// for given timeUnit and timeLength
func (m *DomainFilesProcessed) FetchDataFromDB(db *DB, timeUnit string, timeLength int) {
	iterator := NewDateIterator(db, timeUnit, timeLength)
	for date := iterator.Next(); date.Valid; date = iterator.Next() {
		records := iterator.GetRecordsForDate("snapshot", m, date.Value)
		m.AddDate(date.Value)
		for obj := records.Next(); obj != nil; obj = records.Next() {
			m.AppendValues(obj, "TotalFilesProcessed", iterator.CurrentOffset)
		}
	}
}

// BackupTime is a metric of the total time a backup took
type BackupTime struct {
	metricData
}

// Init initializes the metric
func (m *BackupTime) Init(timeUnit string) {
	m.supportsDomains = false
	m.DateLayout = getDateLayoutForTimeUnit(timeUnit)
	m.Name = "backup time"
	m.Formatter = &TimeFormatter{}
}

// FetchDataFromDB fetches appropriate data from DB and appends values
// for given timeUnit and timeLength
func (m *BackupTime) FetchDataFromDB(db *DB, timeUnit string, timeLength int) {
	iterator := NewDateIterator(db, timeUnit, timeLength)
	for date := iterator.Next(); date.Valid; date = iterator.Next() {
		records := iterator.GetRecordsForDate("repo_backup_times", m, date.Value)
		m.AddDate(date.Value)
		for obj := records.Next(); obj != nil; obj = records.Next() {
			m.AppendValues(obj, "Took", iterator.CurrentOffset)
		}
	}
}

// RepoDiskSpace is a metric of how much space a repo takes on disk
type RepoDiskSpace struct {
	metricData
}

// Init initializes the metric
func (m *RepoDiskSpace) Init(timeUnit string) {
	m.supportsDomains = false
	m.DateLayout = getDateLayoutForTimeUnit(timeUnit)
	m.Name = "disk space occupied"
	m.Formatter = &BytesFormatter{}
}

// SetTitle sets the title of a RepoDiskSpace metric m from given input data
func (m *RepoDiskSpace) SetTitle(name string, repo string, group string, domain string, aggregator string) {
	var d string
	if domain == AllDomains {
		d = "all domains"
	} else {
		d = domain
	}
	m.Title = fmt.Sprintf(
		"%s %s by %s in group %s of repo %s",
		aggregator,
		m.Name,
		d,
		group,
		repo,
	)
}

// FetchDataFromDB fetches appropriate data from DB and appends values
// for given timeUnit and timeLength
func (m *RepoDiskSpace) FetchDataFromDB(db *DB, timeUnit string, timeLength int) {
	iterator := NewDateIterator(db, timeUnit, timeLength)
	for date := iterator.Next(); date.Valid; date = iterator.Next() {
		records := iterator.GetRecordsForDate("repo_raw_data", m, date.Value)
		m.AddDate(date.Value)
		for obj := records.Next(); obj != nil; obj = records.Next() {
			m.AppendValues(obj, "TotalSize", iterator.CurrentOffset)
		}
	}
}

// DomainDiskSpace is a metric of how much space all the backups of a domain takes on disk
type DomainDiskSpace struct {
	metricData
}

// Init initializes the metric
func (m *DomainDiskSpace) Init(timeUnit string) {
	m.supportsDomains = true
	m.DateLayout = getDateLayoutForTimeUnit(timeUnit)
	m.Name = "disk space occupied"
	m.Formatter = &BytesFormatter{}
}

// SetTitle sets the title of a DomainDiskSpace metric m from given input data
func (m *DomainDiskSpace) SetTitle(name string, repo string, group string, domain string, aggregator string) {
	var d string
	if domain == AllDomains {
		d = "all domains"
	} else {
		d = domain
	}
	m.Title = fmt.Sprintf(
		"%s %s by %s in group %s of repo %s",
		aggregator,
		m.Name,
		d,
		group,
		repo,
	)
}

// FetchDataFromDB fetches appropriate data from DB and appends values
// for given timeUnit and timeLength
func (m *DomainDiskSpace) FetchDataFromDB(db *DB, timeUnit string, timeLength int) {
	iterator := NewDateIterator(db, timeUnit, timeLength)
	for date := iterator.Next(); date.Valid; date = iterator.Next() {
		records := iterator.GetRecordsForDate("domain_raw_data", m, date.Value)
		m.AddDate(date.Value)
		for obj := records.Next(); obj != nil; obj = records.Next() {
			m.AppendValues(obj, "TotalSize", iterator.CurrentOffset)
		}
	}
}

// DomainDiskSpaceOnRestore is a metric of how much disk space a domain would take on restore
type DomainDiskSpaceOnRestore struct {
	metricData
}

// Init initializes the metric
func (m *DomainDiskSpaceOnRestore) Init(timeUnit string) {
	m.supportsDomains = true
	m.DateLayout = getDateLayoutForTimeUnit(timeUnit)
	m.Name = "size of restore"
	m.Formatter = &BytesFormatter{}
}

// FetchDataFromDB fetches appropriate data from DB and appends values
// for given timeUnit and timeLength
func (m *DomainDiskSpaceOnRestore) FetchDataFromDB(db *DB, timeUnit string, timeLength int) {
	iterator := NewDateIterator(db, timeUnit, timeLength)
	for date := iterator.Next(); date.Valid; date = iterator.Next() {
		records := iterator.GetRecordsForDate("domain_raw_data", m, date.Value)
		m.AddDate(date.Value)
		for obj := records.Next(); obj != nil; obj = records.Next() {
			m.AppendValues(obj, "TotalSize", iterator.CurrentOffset)
		}
	}
}
