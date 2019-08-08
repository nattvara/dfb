package progress

import (
	"fmt"

	"github.com/nattvara/dfb/internal/restic"
)

// ListenForMessages will listen for restic messages on given channel for report
func (report *Report) ListenForMessages(channel chan restic.Message) {
	for {
		msg := <-channel
		switch msg.Type {
		case "dfb":
			dfb := restic.DFBMessageFromString(msg.Body)
			report.handleDFBMessage(dfb)
		case "status":
			status := restic.StatusMessageFromString(msg.Body)
			report.handleStatusMessage(status)
		case "summary":
			summary := restic.SummaryMessageFromString(msg.Body)
			report.handleSummaryMessage(summary)
		}
	}
}

func (report *Report) handleStatusMessage(msg restic.StatusMessage) {
	if report.currentDomain.Completed {
		return
	}
	report.currentDomain.Elapsed.SetText("Elapsed: " + msg.GetElapsedTime())
	report.currentDomain.ETA.SetText("ETA: " + msg.GetETA())
	report.currentDomain.Files.SetText(fmt.Sprintf("Files: %v/%v", msg.FilesDone, msg.TotalFiles))
	report.currentDomain.Data.SetText(fmt.Sprintf("%s/%s", msg.GetBytesDoneString(), msg.GetTotalBytesString()))
	report.currentDomain.ProgressBar.SetValue(msg.GetProcent())

	if len(msg.CurrentFiles) == 1 {
		report.currentDomain.StatusLines[0].SetText(truncateString(msg.CurrentFiles[0], 100))
		report.currentDomain.StatusLines[1].SetText("")
	} else if len(msg.CurrentFiles) == 2 {
		report.currentDomain.StatusLines[0].SetText(truncateString(msg.CurrentFiles[0], 100))
		report.currentDomain.StatusLines[1].SetText(truncateString(msg.CurrentFiles[1], 100))
	} else {
		report.currentDomain.StatusLines[0].SetText("")
		report.currentDomain.StatusLines[1].SetText("")
	}
}

func truncateString(str string, max int) string {
	out := str
	if len(str) > max {
		if max > 3 {
			max -= 3
		}
		out = str[0:max] + "..."
	}
	return out
}

func (report *Report) handleSummaryMessage(msg restic.SummaryMessage) {
	report.currentDomain.Completed = true
	report.currentDomain.Elapsed.SetText("Took: " + msg.GetDurationString())
	report.currentDomain.Data.SetText("Processed: " + msg.GetDataProcessedString())
	report.currentDomain.Files.SetText("Added: " + msg.GetDataAddedString())
	report.currentDomain.ProgressBar.SetValue(100)
	report.currentDomain.StatusLines[0].SetText("")
	report.currentDomain.StatusLines[1].SetText("")
}

func (report *Report) handleDFBMessage(msg restic.DFBMessage) {
	switch msg.Action {
	case "begin":
		if report.currentDomain != nil {
			report.currentDomain.Completed = true
		}
		report.StartNewDomain(msg.Group, msg.Domain)
	case "unavailable":
		if report.currentDomain != nil {
			report.currentDomain.Completed = true
		}
		report.StartNewEmptyDomain(msg.Group, msg.Domain)
		report.currentDomain.Elapsed.SetText("Took: 0 s")
		report.currentDomain.Data.SetText("Unavailable")
		report.currentDomain.Files.SetText("Added: 0 B")
	case "gathering_stats":
		if report.currentDomain == nil {
			return
		}
		report.currentDomain.Completed = true
		report.currentDomain.StatusLines[0].SetText("gathering stats for " + msg.Domain)
		report.updateLayout()
	case "gathering_stats_done":
		if report.currentDomain == nil {
			return
		}
		report.currentDomain.Completed = true
		report.currentDomain.StatusLines[1].SetText("done.")
		report.updateLayout()
	case "done":
		if report.currentDomain == nil {
			return
		}
		report.done = true
		report.updateLayout()
	case "not_this_repo":
		if report.currentDomain == nil {
			return
		}
		report.currentDomain.Elapsed.SetText("Took: 0 s")
		report.currentDomain.Data.SetText("Skipped")
		report.currentDomain.Files.SetText("Added: 0 B")
		report.updateLayout()
	}
}
