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
	BytesDone        int      `json:"bytes_done"`
	TotalBytes       int      `json:"total_bytes"`
	FilesDone        int      `json:"files_done"`
	TotalFiles       int      `json:"total_files"`
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

// GetElapsedTime returns the elapsed time backup command has been running
func (msg *StatusMessage) GetElapsedTime() string {
	return timeToString(msg.SecondsElapsed)
}

// GetETA returns the estimated time left for backup command
func (msg *StatusMessage) GetETA() string {
	return timeToString(msg.SecondsRemaining)
}

// GetTotalBytesString returns a nicely formatted string of the total bytes to
// process for backup eg. 289.0 MiB 3.1 GiB
func (msg *StatusMessage) GetTotalBytesString() string {
	return bytesToString(msg.TotalBytes)
}

// GetBytesDoneString returns a nicely formatted string of the bytes of data that
// have been processed during backup eg. 289.0 MiB 3.1 GiB
func (msg *StatusMessage) GetBytesDoneString() string {
	return bytesToString(msg.BytesDone)
}

// GetStatusString returns a string representation of status message
func (msg *StatusMessage) GetStatusString() string {
	return fmt.Sprintf(
		"Elapsed: %s, ETA: %s, Files %v/%v Left To Process %s/%s",
		msg.GetElapsedTime(),
		msg.GetETA(),
		msg.FilesDone,
		msg.TotalFiles,
		msg.GetTotalBytesString(),
		msg.GetBytesDoneString(),
	)
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
	return timeToString(int(msg.TotalDuration))
}

// GetDataAddedString returns a nicely formatted string of the amount added,
// eg. 289.0 MiB 3.1 GiB
func (msg *SummaryMessage) GetDataAddedString() string {
	return bytesToString(msg.DataAdded)
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

// timeToString formats a number of seconds s to a shorter representation
// such as 12 s, 31 min or 2.3 h
func timeToString(s int) string {
	var unit string
	var value int
	if s < 60 {
		value = s
		unit = "s"
	} else if s < 60*60 {
		value = s / 60
		unit = "min"
	} else {
		value = s / 60 * 60
		unit = "h"
	}
	return fmt.Sprintf("%v%s", value, unit)
}

// bytesToString formats a float64 of bytes to a nicely formatted
// string eg. 289.0 MiB 3.1 GiB
//
// source: https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func bytesToString(value int) string {
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
