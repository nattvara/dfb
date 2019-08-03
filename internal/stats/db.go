package stats

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-memdb"
	"github.com/nattvara/dfb/internal/paths"
)

// SnapshotSummary is based on a message emitted by restic
// when the backup command finishes
type SnapshotSummary struct {
	// Metadata
	SnapshotID string
	Group      string
	Domain     string
	Repo       string

	// Additional fields used for querying
	GroupWithWildcard  []string
	DomainWithWildcard []string
	DateString         string
	MonthString        string
	YearString         string

	// Values used for metrics
	FilesNew            int
	FilesChanged        int
	FilesUnmodified     int
	DirsNew             int
	DirsChanged         int
	DirsUnmodified      int
	DataBlobs           int
	TreeBlobs           int
	DataAdded           int
	TotalFilesProcessed int
	TotalBytesProcessed int
	TotalDuration       float64
}

type RepoBackupTime struct {
	// Metadata
	ID    int
	Group string
	Repo  string

	// Additional fields used for querying
	GroupWithWildcard []string
	DateString        string
	MonthString       string
	YearString        string

	// Values used for metrics
	Took float64
}

// DB is a database wrapper
//
// Leveraging hashicorp/go-memdb it provides features to load
// memdb with backup data from csv files for a given group,
// and retrieve object by querying various indices
type DB struct {
	memdb *memdb.MemDB
}

// Load loads db with data from csv files for given group
func (db *DB) Load(groupName string) {
	statsDir := fmt.Sprintf("%s/%s/stats", paths.DFB(), groupName)

	for _, record := range csvReadSummaries(statsDir + "/snapshots.csv") {
		db.InsertRecord("snapshot", record)
	}
	for _, record := range csvReadRepoBackupTime(statsDir + "/repo_time_took.csv") {
		db.InsertRecord("repo_backup_times", record)
	}
}

// InsertRecord will insert a record in given table in memdb instance
func (db *DB) InsertRecord(table string, record interface{}) {
	txn := db.memdb.Txn(true)
	if err := txn.Insert(table, record); err != nil {
		panic(err)
	}
	txn.Commit()
}

// GetIndexFromTimeUnit returns the index to use for given time unit, use includeDomain
// to search for a specifc domain (not supported by all indices)
func (db *DB) GetIndexFromTimeUnit(timeUnit string, includeDomain bool) string {
	index := "repo_group"
	if includeDomain {
		index += "_domain"
	}
	switch strings.ToLower(timeUnit) {
	case TimeUnitDays:
		index += "_daily"
	case TimeUnitMonths:
		index += "_monthly"
	case TimeUnitYears:
		index += "_yearly"
	default:
		log.Fatal("unsupported timeUnit: " + timeUnit)
	}
	return index
}

// NewDB returns a new DB instance
func NewDB() *DB {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"snapshot": &memdb.TableSchema{
				Name: "snapshot",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "SnapshotID"},
					},
					"repo_group_domain_daily": &memdb.IndexSchema{
						Name:   "repo_group_domain_daily",
						Unique: false,
						Indexer: &memdb.CompoundMultiIndex{
							AllowMissing: false,
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Repo"},
								&memdb.StringSliceFieldIndex{Field: "GroupWithWildcard"},
								&memdb.StringSliceFieldIndex{Field: "DomainWithWildcard"},
								&memdb.StringFieldIndex{Field: "DateString"},
							},
						},
					},
					"repo_group_domain_monthly": &memdb.IndexSchema{
						Name:   "repo_group_domain_monthly",
						Unique: false,
						Indexer: &memdb.CompoundMultiIndex{
							AllowMissing: false,
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Repo"},
								&memdb.StringSliceFieldIndex{Field: "GroupWithWildcard"},
								&memdb.StringSliceFieldIndex{Field: "DomainWithWildcard"},
								&memdb.StringFieldIndex{Field: "MonthString"},
							},
						},
					},
					"repo_group_domain_yearly": &memdb.IndexSchema{
						Name:   "repo_group_domain_yearly",
						Unique: false,
						Indexer: &memdb.CompoundMultiIndex{
							AllowMissing: false,
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Repo"},
								&memdb.StringSliceFieldIndex{Field: "GroupWithWildcard"},
								&memdb.StringSliceFieldIndex{Field: "DomainWithWildcard"},
								&memdb.StringFieldIndex{Field: "YearString"},
							},
						},
					},
				},
			},
			"repo_backup_times": &memdb.TableSchema{
				Name: "repo_backup_times",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "ID"},
					},
					"repo_group_daily": &memdb.IndexSchema{
						Name:   "repo_group_daily",
						Unique: false,
						Indexer: &memdb.CompoundMultiIndex{
							AllowMissing: false,
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Repo"},
								&memdb.StringFieldIndex{Field: "Group"},
								&memdb.StringFieldIndex{Field: "DateString"},
							},
						},
					},
					"repo_group_monthly": &memdb.IndexSchema{
						Name:   "repo_group_monthly",
						Unique: false,
						Indexer: &memdb.CompoundMultiIndex{
							AllowMissing: false,
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Repo"},
								&memdb.StringFieldIndex{Field: "Group"},
								&memdb.StringFieldIndex{Field: "MonthString"},
							},
						},
					},
					"repo_group_yearly": &memdb.IndexSchema{
						Name:   "repo_group_yearly",
						Unique: false,
						Indexer: &memdb.CompoundMultiIndex{
							AllowMissing: false,
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{Field: "Repo"},
								&memdb.StringFieldIndex{Field: "Group"},
								&memdb.StringFieldIndex{Field: "YearString"},
							},
						},
					},
				},
			},
		},
	}

	memdb, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	return &DB{
		memdb: memdb,
	}
}
