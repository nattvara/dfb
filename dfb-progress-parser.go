package main

import (
	"bufio"
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
	PercentDone float64 `json:"percent_done"`
}

// GetProcenString returns a string with ProcentDone formatted like X%
func (msg *StatusMessage) GetProcenString() string {
	return fmt.Sprintf("%v%%", math.Round(msg.PercentDone*100))
}

// SummaryMessage should match summary messages from restic which
// are emitted once a backup command is completed
type SummaryMessage struct {
	FilesNew      int     `json:"files_new"`
	FilesChanged  int     `json:"files_changed"`
	DirsNew       int     `json:"dirs_new"`
	DirsChanged   int     `json:"dirs_changed"`
	DataAadded    int     `json:"data_added"`
	TotalDuration float64 `json:"total_duration"`
}

// GetDurationString returns a string with TotalDuration formatted like x.xs
func (msg *SummaryMessage) GetDurationString() string {
	return fmt.Sprintf("%vs", math.Round(msg.TotalDuration*10)/10)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: dfb-progress-parser [output-prefix]")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	prefix := os.Args[1]

	for scanner.Scan() {
		var msg Message
		json.Unmarshal(scanner.Bytes(), &msg)

		switch msg.Type {
		case "":
		case "status":
			var status StatusMessage
			json.Unmarshal(scanner.Bytes(), &status)

			printStatusMessage(&status, prefix)

		case "summary":
			var summary SummaryMessage
			json.Unmarshal(scanner.Bytes(), &summary)

			printSummaryMessage(&summary, prefix)
		default:
			printUnknownMessage(&msg)
		}
	}
}

func printStatusMessage(msg *StatusMessage, prefix string) {
	clearLine()
	fmt.Printf("%s \t %s", prefix, msg.GetProcenString())
}

func printSummaryMessage(msg *SummaryMessage, prefix string) {
	clearLine()
	fmt.Printf("%s \t 100%% ⏱  %v \n", prefix, msg.GetDurationString())
}

func printUnknownMessage(msg *Message) {
	fmt.Printf("unkown message type %s\n", msg.Type)
}

func clearLine() {
	runes := make([]rune, tm.Width())
	fmt.Print(string(runes))
	fmt.Print("\r \r")
}
