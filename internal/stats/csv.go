package stats

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// csvReadSummaries reads given csv file, parses and returns snapshot summaries
func csvReadSummaries(filename string) []*SnapshotSummary {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var summaries []*SnapshotSummary

	reader := csv.NewReader(f)
	lineNumber := -1
	for {
		line, err := reader.Read()
		lineNumber++
		if err == io.EOF {
			break
		}

		if err != nil {
			if err, ok := err.(*csv.ParseError); ok && err.Err == csv.ErrFieldCount {
				fmt.Printf("delete line %v\n", lineNumber)
				continue
			}
		}

		filesNew, _ := strconv.Atoi(line[1])
		filesChanged, _ := strconv.Atoi(line[2])
		filesUnmodified, _ := strconv.Atoi(line[3])
		dirsNew, _ := strconv.Atoi(line[4])
		dirsChanged, _ := strconv.Atoi(line[5])
		dirsUnmodified, _ := strconv.Atoi(line[6])
		dataBlobs, _ := strconv.Atoi(line[7])
		treeBlobs, _ := strconv.Atoi(line[8])
		dataAdded, _ := strconv.Atoi(line[9])
		totalFilesProcessed, _ := strconv.Atoi(line[10])
		totalBytesProcessed, _ := strconv.Atoi(line[11])
		totalDuration, _ := strconv.ParseFloat(line[12], 64)
		date, _ := time.Parse("2006-01-02T15:04:05Z0700", line[17])

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

			SnapshotID: line[13],
			Group:      line[14],
			Domain:     line[15],
			Repo:       line[16],

			GroupWithWildcard:  []string{line[14], AllDomains},
			DomainWithWildcard: []string{line[15], AllDomains},
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
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	var backupTimes []*RepoBackupTime

	reader := csv.NewReader(f)
	lineNumber := -1
	for {
		line, err := reader.Read()
		lineNumber++
		if err == io.EOF {
			break
		}

		if err != nil {
			if err, ok := err.(*csv.ParseError); ok && err.Err == csv.ErrFieldCount {
				log.Printf("%s parse error at line %v\n", filename, lineNumber)
				continue
			}
		}

		duration, _ := strconv.ParseFloat(line[0], 64)
		date, _ := time.Parse("2006-01-02T15:04:05Z0700", line[3])

		bt := &RepoBackupTime{
			Took: duration,

			ID:    lineNumber,
			Group: line[1],
			Repo:  line[2],

			DateString:  date.Format(getDateLayoutForTimeUnit(TimeUnitDays)),
			MonthString: date.Format(getDateLayoutForTimeUnit(TimeUnitMonths)),
			YearString:  date.Format(getDateLayoutForTimeUnit(TimeUnitYears)),
		}
		backupTimes = append(backupTimes, bt)
	}

	return backupTimes
}
