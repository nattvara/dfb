package stats

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-memdb"
	"github.com/nattvara/dfb/src/internal/paths"
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

// DB is a database wrapper
//
// Leveraging hashicorp/go-memdb it provides features to load
// memdb with backup data from csv files for a given group,
// and retrieve object by querying various indices
type DB struct {
	memdb *memdb.MemDB
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

// Load loads db with data from csv files for given group
func (db *DB) Load(groupName string) {
	statsDir := fmt.Sprintf("%s/%s/stats", paths.DFB(), groupName)

	summaries := csvReadSummaries(statsDir + "/snapshots.csv")
	db.InsertSnapshotSummaries(summaries)
}

// InsertSnapshotSummaries inserts memdb instance with snapshot summaries
func (db *DB) InsertSnapshotSummaries(summaries []*SnapshotSummary) {
	txn := db.memdb.Txn(true)
	for _, snapshot := range summaries {
		if err := txn.Insert("snapshot", snapshot); err != nil {
			panic(err)
		}
	}
	txn.Commit()
}

// GetIndexFromTimeUnit returns the index to use for given time unit
func (db *DB) GetIndexFromTimeUnit(timeUnit string) string {
	var index string
	switch strings.ToLower(timeUnit) {
	case TimeUnitDays:
		index = "repo_group_domain_daily"
	case TimeUnitMonths:
		index = "repo_group_domain_monthly"
	case TimeUnitYears:
		index = "repo_group_domain_yearly"
	default:
		log.Fatal("unsupported timeUnit: " + timeUnit)
	}
	return index
}
