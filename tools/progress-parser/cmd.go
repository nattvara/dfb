package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	d "github.com/nattvara/dfb/internal/domains"
	"github.com/nattvara/dfb/internal/paths"
	"github.com/nattvara/dfb/internal/restic"

	tm "github.com/buger/goterm"
)

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

		var msg restic.Message
		json.Unmarshal(scanner.Bytes(), &msg)

		switch msg.Type {
		case "":
		case "status":
			var status restic.StatusMessage
			json.Unmarshal(scanner.Bytes(), &status)
			linesPrinted = PrintStatusMessage(status, domain)
		case "summary":
			var summary restic.SummaryMessage
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
		tm.Flush()
	}
	tm.Print(tm.RESET_LINE)
	tm.Flush()
}

// PrintStatusMessage prints a status message with procent done, ETA and which files are currently being backed up
func PrintStatusMessage(msg restic.StatusMessage, domain d.Domain) int {
	var linesPrinted int

	message := fmt.Sprintf("  backing up %s", domain.Name)
	tm.Printf(message)

	tm.MoveCursorForward(50 - len(message))
	tm.Print(msg.GetProcentString())
	if msg.SecondsRemaining != 0 {
		tm.Print(" " + msg.GetStatusString())
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
func PrintSummaryMessage(msg restic.SummaryMessage, domain d.Domain) int {
	var linesPrinted int

	message := fmt.Sprintf("  backing up %s", domain.Name)
	tm.Printf(message)
	tm.MoveCursorForward(50 - len(message))
	tm.Printf("100%% ⏱  %s 💾 %s 📊 gathering stats... ", msg.GetDurationString(), msg.GetDataAddedString())

	tm.Flush()
	return linesPrinted
}
