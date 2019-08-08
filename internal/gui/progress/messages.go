package progress

import (
	"fmt"

	"github.com/nattvara/dfb/internal/restic"
)

// ListenForMessages will listen for restic messages on given channel for gui
func (gui *ProgressGUI) ListenForMessages(channel chan restic.Message) {
	for {
		msg := <-channel
		switch msg.Type {
		case "dfb":
			dfb := restic.DFBMessageFromString(msg.Body)
			gui.handleDFBMessage(dfb)
		case "status":
			status := restic.StatusMessageFromString(msg.Body)
			gui.handleStatusMessage(status)
		case "summary":
			summary := restic.SummaryMessageFromString(msg.Body)
			gui.handleSummaryMessage(summary)
		}
	}
}

func (gui *ProgressGUI) handleStatusMessage(msg restic.StatusMessage) {
	if gui.currentDomain.Completed {
		return
	}
	gui.currentDomain.Elapsed.SetText("Elapsed: " + msg.GetElapsedTime())
	gui.currentDomain.ETA.SetText("ETA: " + msg.GetETA())
	gui.currentDomain.Files.SetText(fmt.Sprintf("Files: %v/%v", msg.FilesDone, msg.TotalFiles))
	gui.currentDomain.Data.SetText(fmt.Sprintf("%s/%s", msg.GetBytesDoneString(), msg.GetTotalBytesString()))
	gui.currentDomain.ProgressBar.SetValue(msg.GetProcent())

	if len(msg.CurrentFiles) == 1 {
		gui.currentDomain.StatusLines[0].SetText(truncateString(msg.CurrentFiles[0], 100))
		gui.currentDomain.StatusLines[1].SetText("")
	} else if len(msg.CurrentFiles) == 2 {
		gui.currentDomain.StatusLines[0].SetText(truncateString(msg.CurrentFiles[0], 100))
		gui.currentDomain.StatusLines[1].SetText(truncateString(msg.CurrentFiles[1], 100))
	} else {
		gui.currentDomain.StatusLines[0].SetText("")
		gui.currentDomain.StatusLines[1].SetText("")
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

func (gui *ProgressGUI) handleSummaryMessage(msg restic.SummaryMessage) {
	gui.currentDomain.Completed = true
	gui.currentDomain.Elapsed.SetText("Took: " + msg.GetDurationString())
	gui.currentDomain.Data.SetText("Processed: " + msg.GetDataProcessedString())
	gui.currentDomain.Files.SetText("Added: " + msg.GetDataAddedString())
	gui.currentDomain.ProgressBar.SetValue(100)
	gui.currentDomain.StatusLines[0].SetText("")
	gui.currentDomain.StatusLines[1].SetText("")
}

func (gui *ProgressGUI) handleDFBMessage(msg restic.DFBMessage) {
	switch msg.Action {
	case "begin":
		if gui.currentDomain != nil {
			gui.currentDomain.Completed = true
		}
		gui.StartNewDomain(msg.Group, msg.Domain)
	case "unavailable":
		if gui.currentDomain != nil {
			gui.currentDomain.Completed = true
		}
		gui.StartNewEmptyDomain(msg.Group, msg.Domain)
		gui.currentDomain.Elapsed.SetText("Took: 0 s")
		gui.currentDomain.Data.SetText("Unavailable")
		gui.currentDomain.Files.SetText("Added: 0 B")
	case "gathering_stats":
		if gui.currentDomain == nil {
			return
		}
		gui.currentDomain.Completed = true
		gui.currentDomain.StatusLines[0].SetText("gathering stats for " + msg.Domain)
		gui.updateLayout()
	case "gathering_stats_done":
		if gui.currentDomain == nil {
			return
		}
		gui.currentDomain.Completed = true
		gui.currentDomain.StatusLines[1].SetText("done.")
		gui.updateLayout()
	case "done":
		if gui.currentDomain == nil {
			return
		}
		gui.Done = true
		gui.updateLayout()
	case "not_this_repo":
		if gui.currentDomain == nil {
			return
		}
		gui.currentDomain.Elapsed.SetText("Took: 0 s")
		gui.currentDomain.Data.SetText("Skipped")
		gui.currentDomain.Files.SetText("Added: 0 B")
		gui.updateLayout()
	}
}
