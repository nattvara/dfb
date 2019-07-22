package main

import (
	"bufio"
	d "dfb/src/internal/domains"
	"dfb/src/internal/paths"
	"encoding/json"
	"fmt"
	"math"
	"os"

	tm "github.com/buger/goterm"
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

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: dfb-progress-parser [group] [domain]")
		os.Exit(1)
	}

	groupName := os.Args[1]
	domainName := os.Args[2]

	domain := d.Load(
		domainName,
		groupName,
		fmt.Sprintf("%s/%s", paths.DFB(), groupName),
	)

	tm.Flush()
	var linesPrinted int

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		ClearPreviousLines(linesPrinted)

		var msg Message
		json.Unmarshal(scanner.Bytes(), &msg)

		switch msg.Type {
		case "":
		case "status":
			var status StatusMessage
			json.Unmarshal(scanner.Bytes(), &status)
			linesPrinted = PrintStatusMessage(status, domain)
		case "summary":
			var summary SummaryMessage
			json.Unmarshal(scanner.Bytes(), &summary)

			linesPrinted = PrintSummaryMessage(summary, domain)
		}
	}
}

// ClearPreviousLines takes the number of lines last printed, clears them at sets
func ClearPreviousLines(number int) {
	for i := 1; i <= number; i++ {
		tm.Print(tm.RESET_LINE)
		tm.MoveCursorUp(1)
	}
	tm.Flush()
}

// PrintStatusMessage prints a status message with procent done, ETA and which files are currently being backed up
func PrintStatusMessage(msg StatusMessage, domain d.Domain) int {
	var linesPrinted int

	message := fmt.Sprintf("  backing up %s", domain.Name)
	tm.Printf(message)

	tm.MoveCursorForward(50 - len(message))
	tm.Print(msg.GetProcenString())
	if msg.SecondsRemaining != 0 {
		tm.Printf("  ETA %vs", msg.SecondsRemaining)
	}

	tm.Printf("\n")
	linesPrinted++

	if len(msg.CurrentFiles) == 1 {
		tm.Println(msg.CurrentFiles[0])
		linesPrinted++
	} else if len(msg.CurrentFiles) == 2 {
		tm.Printf("  current files %s\n", msg.CurrentFiles[0])
		tm.Printf("  current files %s\n", msg.CurrentFiles[1])
		linesPrinted += 2
	}

	tm.Flush()
	return linesPrinted
}

// PrintSummaryMessage prints a summary message with time taken to perform backup of domain
func PrintSummaryMessage(msg SummaryMessage, domain d.Domain) int {
	var linesPrinted int

	message := fmt.Sprintf("  backing up %s", domain.Name)
	tm.Printf(message)
	tm.MoveCursorForward(50 - len(message))
	tm.Printf("100%% ⏱  %s \n", msg.GetDurationString())

	tm.Flush()
	return linesPrinted
}
