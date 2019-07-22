package restic

import (
	"fmt"
	"math"
)

// Message should match any restic message
type Message struct {
	Type string `json:"message_type"`
}

// StatusMessage should match status messages from restic
// which are emitted during backup
type StatusMessage struct {
	PercentDone      float64  `json:"percent_done"`
	SecondsElapsed   int      `json:"seconds_elapsed"`
	SecondsRemaining int      `json:"seconds_remaining"`
	CurrentFiles     []string `json:"current_files"`
}

// GetProcenString returns a string with ProcentDone formatted like X%
func (msg *StatusMessage) GetProcenString() string {
	return fmt.Sprintf("%v%%", math.Round(msg.PercentDone*100))
}

// SummaryMessage should match summary messages from restic which
// are emitted once a backup command is completed
type SummaryMessage struct {
	FilesNew       int     `json:"files_new"`
	FilesChanged   int     `json:"files_changed"`
	FilesProcessed int     `json:"total_files_processed"`
	DirsNew        int     `json:"dirs_new"`
	DirsChanged    int     `json:"dirs_changed"`
	BytesProcessed int     `json:"total_bytes_processed"`
	DataAdded      int     `json:"data_added"`
	TotalDuration  float64 `json:"total_duration"`
	SnapshotID     string  `json:"snapshot_id"`
}

// GetDurationString returns a string with TotalDuration formatted like x.xs
func (msg *SummaryMessage) GetDurationString() string {
	return fmt.Sprintf("%vs", math.Round(msg.TotalDuration*10)/10)
}
