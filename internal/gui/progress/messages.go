package progress

import (
	"github.com/nattvara/dfb/internal/restic"
)

// MessageReceiver receives json messages and updates Report accordingly
type MessageReceiver struct {
	Report *Report
}

// ListenForMessages will listen for restic messages on given channel
func (receiver *MessageReceiver) ListenForMessages(channel chan restic.Message) {
	for {
		msg := <-channel
		switch msg.Type {
		case "dfb":
			dfb := restic.DFBMessageFromString(msg.Body)
			receiver.handleDFBMessage(dfb)
		case "status":
			status := restic.StatusMessageFromString(msg.Body)
			receiver.handleStatusMessage(status)
		case "summary":
			summary := restic.SummaryMessageFromString(msg.Body)
			receiver.handleSummaryMessage(summary)
		}
	}
}

func (receiver *MessageReceiver) handleStatusMessage(msg restic.StatusMessage) {
	receiver.Report.StatusComponent.SetElapsed(msg.GetElapsedTime())
	receiver.Report.StatusComponent.SetETA(msg.GetETA())
	receiver.Report.StatusComponent.SetFilesDone(msg.FilesDone, msg.TotalFiles)
	receiver.Report.StatusComponent.SetBytesDone(msg.GetBytesDoneString(), msg.GetTotalBytesString())
	receiver.Report.StatusComponent.SetProgress(msg.GetProcent())
	receiver.Report.StatusComponent.SetStatusLinesFromCurrentFiles(msg.CurrentFiles)
}

func (receiver *MessageReceiver) handleSummaryMessage(msg restic.SummaryMessage) {
	receiver.Report.StatusComponent.SetProgress(100)
	receiver.Report.CompleteCurrentDomain(msg)
}

func (receiver *MessageReceiver) handleDFBMessage(msg restic.DFBMessage) {
	switch msg.Action {
	case "begin":
		receiver.Report.StartNewDomain(msg.Group, msg.Domain)
	case "unavailable":
		receiver.Report.StartNewDomain(msg.Group, msg.Domain)
		receiver.Report.CompleteUnavailibleDomain(msg.Group, msg.Domain, "Unavailable")
	case "gathering_stats":
		receiver.Report.StatusComponent.SetFirstStatusLine("gathering stats for " + msg.Domain)
		receiver.Report.StatusComponent.SetSecondStatusLine("")
	case "gathering_stats_done":
		receiver.Report.StatusComponent.SetSecondStatusLine("done.")
	case "not_this_repo":
		receiver.Report.CompleteUnavailibleDomain(msg.Group, msg.Domain, "Not this repo")
	case "done":
		receiver.Report.Done()
	}
}
