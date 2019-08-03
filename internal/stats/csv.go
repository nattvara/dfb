package stats

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// csvFileIterator is a type that can iterate over the records in a csv file
type csvFileIterator struct {
	filename          string
	file              *os.File
	reader            *csv.Reader
	CurrentLineNumber int
}

// Open opens csv file at given filename for csvFileIterator it
func (it *csvFileIterator) Open(filename string) {
	it.filename = filename

	var err error
	it.file, err = os.Open(filename)

	if err != nil {
		panic("failed to open csv file. " + err.Error())
	}

	it.reader = csv.NewReader(it.file)
	it.CurrentLineNumber = -1
}

// Close closes file descriptor used for reading csv by csvFileIterator it
func (it *csvFileIterator) Close() {
	it.file.Close()
}

// Next reads and returns the next record from opened csv file by csvFileIterator it
func (it *csvFileIterator) Next() []string {
	it.CurrentLineNumber++
	line, err := it.reader.Read()
	if err == io.EOF {
		return nil
	}

	if err != nil {
		if err, ok := err.(*csv.ParseError); ok && err.Err == csv.ErrFieldCount {
			fmt.Printf("%s parse error at line %v\n", it.filename, it.CurrentLineNumber)
			return it.Next()
		}
	}

	return line
}

// csvReadSummaries reads given csv file, parses and returns snapshot summaries
func csvReadSummaries(filename string) []*SnapshotSummary {
	var summaries []*SnapshotSummary

	it := &csvFileIterator{}
	it.Open(filename)
	defer it.Close()

	for record := it.Next(); record != nil; record = it.Next() {
		filesNew, _ := strconv.Atoi(record[1])
		filesChanged, _ := strconv.Atoi(record[2])
		filesUnmodified, _ := strconv.Atoi(record[3])
		dirsNew, _ := strconv.Atoi(record[4])
		dirsChanged, _ := strconv.Atoi(record[5])
		dirsUnmodified, _ := strconv.Atoi(record[6])
		dataBlobs, _ := strconv.Atoi(record[7])
		treeBlobs, _ := strconv.Atoi(record[8])
		dataAdded, _ := strconv.Atoi(record[9])
		totalFilesProcessed, _ := strconv.Atoi(record[10])
		totalBytesProcessed, _ := strconv.Atoi(record[11])
		totalDuration, _ := strconv.ParseFloat(record[12], 64)
		date, _ := time.Parse("2006-01-02T15:04:05Z0700", record[17])

		summary := &SnapshotSummary{
			FilesNew:            filesNew,
			FilesChanged:        filesChanged,
			FilesUnmodified:     filesUnmodified,
			DirsNew:             dirsNew,
			DirsChanged:         dirsChanged,
			DirsUnmodified:      dirsUnmodified,
			DataBlobs:           dataBlobs,
			TreeBlobs:           treeBlobs,
			DataAdded:           dataAdded,
			TotalFilesProcessed: totalFilesProcessed,
			TotalBytesProcessed: totalBytesProcessed,
			TotalDuration:       totalDuration,

			SnapshotID: record[13],
			Group:      record[14],
			Domain:     record[15],
			Repo:       record[16],

			GroupWithWildcard:  []string{record[14], AllDomains},
			DomainWithWildcard: []string{record[15], AllDomains},
			DateString:         date.Format(getDateLayoutForTimeUnit(TimeUnitDays)),
			MonthString:        date.Format(getDateLayoutForTimeUnit(TimeUnitMonths)),
			YearString:         date.Format(getDateLayoutForTimeUnit(TimeUnitYears)),
		}
		summaries = append(summaries, summary)
	}

	return summaries
}

// csvReadRepoBackupTime reads given csv file, parses and returns repo backup times
func csvReadRepoBackupTime(filename string) []*RepoBackupTime {
	var backupTimes []*RepoBackupTime

	it := &csvFileIterator{}
	it.Open(filename)
	defer it.Close()

	for record := it.Next(); record != nil; record = it.Next() {

		duration, _ := strconv.ParseFloat(record[0], 64)
		date, _ := time.Parse("2006-01-02T15:04:05Z0700", record[3])

		bt := &RepoBackupTime{
			Took: duration,

			ID:    it.CurrentLineNumber,
			Group: record[1],
			Repo:  record[2],

			DateString:  date.Format(getDateLayoutForTimeUnit(TimeUnitDays)),
			MonthString: date.Format(getDateLayoutForTimeUnit(TimeUnitMonths)),
			YearString:  date.Format(getDateLayoutForTimeUnit(TimeUnitYears)),
		}
		backupTimes = append(backupTimes, bt)
	}

	return backupTimes
}

// csvReadRepoRawData reads given csv file, parses and returns repo raw data records
func csvReadRepoRawData(filename string) []*RepoRawData {
	var rawData []*RepoRawData

	it := &csvFileIterator{}
	it.Open(filename)
	defer it.Close()

	for record := it.Next(); record != nil; record = it.Next() {

		totalSize, _ := strconv.ParseInt(record[0], 10, 64)
		totalFileCount, _ := strconv.Atoi(record[1])
		totalBlobCount, _ := strconv.Atoi(record[2])
		date, _ := time.Parse("2006-01-02T15:04:05Z0700", record[5])

		rd := &RepoRawData{
			TotalSize:      totalSize,
			TotalFileCount: totalFileCount,
			TotalBlobCount: totalBlobCount,

			ID:    it.CurrentLineNumber,
			Group: record[3],
			Repo:  record[4],

			DateString:  date.Format(getDateLayoutForTimeUnit(TimeUnitDays)),
			MonthString: date.Format(getDateLayoutForTimeUnit(TimeUnitMonths)),
			YearString:  date.Format(getDateLayoutForTimeUnit(TimeUnitYears)),
		}
		rawData = append(rawData, rd)
	}

	return rawData
}

// csvReadDomainRawData reads given csv file, parses and returns domain raw data records
func csvReadDomainRawData(filename string) []*DomainRawData {
	var rawData []*DomainRawData

	it := &csvFileIterator{}
	it.Open(filename)
	defer it.Close()

	for record := it.Next(); record != nil; record = it.Next() {

		totalSize, _ := strconv.ParseInt(record[0], 10, 64)
		totalFileCount, _ := strconv.Atoi(record[1])
		totalBlobCount, _ := strconv.Atoi(record[2])
		date, _ := time.Parse("2006-01-02T15:04:05Z0700", record[6])

		rd := &DomainRawData{
			TotalSize:      totalSize,
			TotalFileCount: totalFileCount,
			TotalBlobCount: totalBlobCount,

			ID:     it.CurrentLineNumber,
			Group:  record[3],
			Domain: record[4],
			Repo:   record[5],

			GroupWithWildcard:  []string{record[3], AllDomains},
			DomainWithWildcard: []string{record[4], AllDomains},

			DateString:  date.Format(getDateLayoutForTimeUnit(TimeUnitDays)),
			MonthString: date.Format(getDateLayoutForTimeUnit(TimeUnitMonths)),
			YearString:  date.Format(getDateLayoutForTimeUnit(TimeUnitYears)),
		}
		rawData = append(rawData, rd)
	}

	return rawData
}

// csvReadDomainRestoreSize reads given csv file, parses and returns domain restore size records
func csvReadDomainRestoreSize(filename string) []*DomainRestoreSize {
	var restoreSizes []*DomainRestoreSize

	it := &csvFileIterator{}
	it.Open(filename)
	defer it.Close()

	for record := it.Next(); record != nil; record = it.Next() {

		totalSize, _ := strconv.ParseInt(record[0], 10, 64)
		totalFileCount, _ := strconv.Atoi(record[1])
		date, _ := time.Parse("2006-01-02T15:04:05Z0700", record[5])

		rs := &DomainRestoreSize{
			TotalSize:      totalSize,
			TotalFileCount: totalFileCount,

			ID:     it.CurrentLineNumber,
			Group:  record[2],
			Domain: record[3],
			Repo:   record[4],

			GroupWithWildcard:  []string{record[2], AllDomains},
			DomainWithWildcard: []string{record[3], AllDomains},

			DateString:  date.Format(getDateLayoutForTimeUnit(TimeUnitDays)),
			MonthString: date.Format(getDateLayoutForTimeUnit(TimeUnitMonths)),
			YearString:  date.Format(getDateLayoutForTimeUnit(TimeUnitYears)),
		}
		restoreSizes = append(restoreSizes, rs)
	}

	return restoreSizes
}
