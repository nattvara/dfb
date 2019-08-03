package restic

import (
	"encoding/json"
	"fmt"
	"math"
)

// Message should match any restic message
type Message struct {
	Type string `json:"message_type"`
	Body string
}

// StatusMessage should match status messages from restic
// which are emitted during backup
type StatusMessage struct {
	PercentDone      float64  `json:"percent_done"`
	SecondsElapsed   int      `json:"seconds_elapsed"`
	SecondsRemaining int      `json:"seconds_remaining"`
	CurrentFiles     []string `json:"current_files"`
}

// GetProcentString returns a string with ProcentDone formatted like X%
func (msg *StatusMessage) GetProcentString() string {
	return fmt.Sprintf("%v%%", math.Round(msg.PercentDone*100))
}

// GetProcent returns a nice rounded value of PercentDone
func (msg *StatusMessage) GetProcent() float64 {
	return math.Round(msg.PercentDone * 100)
}

// GetETAString will return a string with time taken and current ETA
func (msg *StatusMessage) GetETAString() string {
	return fmt.Sprintf("time: %vs, ETA: %vs", msg.SecondsElapsed, msg.SecondsRemaining)
}

// StatusMessageFromString will create a StatusMessage from given string
func StatusMessageFromString(data string) StatusMessage {
	var status StatusMessage
	json.Unmarshal([]byte(data), &status)
	return status
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

// GetDataAddedString returns a nicely formatted string of the amount added,
// eg. 289.0 MiB 3.1 GiB
// source: https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func (msg *SummaryMessage) GetDataAddedString() string {
	bytes := msg.DataAdded
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
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

// SummaryMessageFromString will create a SummaryMessage from given string
func SummaryMessageFromString(data string) SummaryMessage {
	var summary SummaryMessage
	json.Unmarshal([]byte(data), &summary)
	return summary
}

// DFBMessage is not a message from restic, but a custom dfb message,
// will be sent in roughly the same context as restic messages (meant to
// be parsed the same)
type DFBMessage struct {
	Group  string `json:"group"`
	Domain string `json:"domain"`
	Action string `json:"action"`
}

// DFBMessageFromString will create a DFBMessage from given string
func DFBMessageFromString(data string) DFBMessage {
	var dfb DFBMessage
	json.Unmarshal([]byte(data), &dfb)
	return dfb
}
