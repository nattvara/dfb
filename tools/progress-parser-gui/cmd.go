// progress-gui will launch a gui window to display the gui of
// a backup.
// progress gui will parse messages from dfb and restic on stdin like the
// following:
//
// {"message_type":"dfb","action":"begin","group":"some-group","domain":"some-domain"}
// {"message_type":"summary","files_new":1,"files_changed":2,"files_unmodified":83,"dirs_new":0,"dirs_changed":0,"dirs_unmodified":0,"data_blobs":0,"tree_blobs":0,"data_added":0,"total_files_processed":83,"total_bytes_processed":43535,"total_duration":0.388768151,"snapshot_id":"xxx"}
//
// gui-progress is based on fyne.io, see:
// https://github.com/fyne-io/fyne

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	d "github.com/nattvara/dfb/internal/domains"
	"github.com/nattvara/dfb/internal/paths"
	"github.com/nattvara/dfb/internal/restic"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

func main() {

	app := app.New()

	progress := NewProgress(app)
	progress.LoadUI(app)

	messages := make(chan restic.Message)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			var msg restic.Message
			json.Unmarshal(scanner.Bytes(), &msg)
			msg.Body = scanner.Text()
			messages <- msg
		}
	}()

	go progress.ListenForMessages(messages)
	app.Run()
}

// DomainProgress contains a domain and the widgets to display the backup progress
// of that domain
type DomainProgress struct {
	Domain    d.Domain
	Completed bool

	Name        *widget.Label
	Elapsed     *widget.Label
	ETA         *widget.Label
	Files       *widget.Label
	Data        *widget.Label
	ProgressBar *widget.ProgressBar
	StatusLines []*widget.Label
}

// NewProgress creates a new gui fyne app
func NewProgress(app fyne.App) *ProgressGUI {
	gui := &ProgressGUI{}
	gui.app = app
	return gui
}

// ProgressGUI contains fyne and app related values
type ProgressGUI struct {
	domains       []*DomainProgress
	currentDomain *DomainProgress
	Done          bool

	window fyne.Window
	app    fyne.App
}

// LoadUI will load the initial UI for gui
func (gui *ProgressGUI) LoadUI(app fyne.App) {
	now := time.Now().Format("15:04")
	gui.window = app.NewWindow("Progress report for dfb backup started at " + now)

	gui.window.SetContent(widget.NewLabel("waiting for messages on stdin"))
	gui.window.Show()
}

func (gui *ProgressGUI) updateLayout() {
	current := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(1000, 200)))
	for _, domain := range gui.domains {
		if domain.Completed {
			continue
		}

		name := fyne.NewContainerWithLayout(
			layout.NewGridLayout(1),
			domain.Name,
		)
		progress := fyne.NewContainerWithLayout(
			layout.NewGridLayout(1),
			domain.ProgressBar,
		)
		stats := fyne.NewContainerWithLayout(
			layout.NewGridLayout(4),
			domain.Elapsed,
			domain.ETA,
			domain.Files,
			domain.Data,
		)
		line1 := fyne.NewContainerWithLayout(
			layout.NewGridLayout(1),
			gui.currentDomain.StatusLines[0],
		)
		line2 := fyne.NewContainerWithLayout(
			layout.NewGridLayout(1),
			gui.currentDomain.StatusLines[1],
		)

		current.AddObject(fyne.NewContainerWithLayout(
			layout.NewGridLayout(1),
			name,
			progress,
			stats,
			line1,
			line2,
		))
	}

	done := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(1000, 80)))
	if gui.Done {
		done.AddObject(fyne.NewContainerWithLayout(
			layout.NewGridLayout(1),
			widget.NewLabelWithStyle("--- Backup Finished ---", fyne.TextAlignCenter, fyne.TextStyle{}),
		))
	}

	completed := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	for _, domain := range gui.domains {
		if !domain.Completed {
			continue
		}
		completed.AddObject(
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(2),
				domain.Name,
				fyne.NewContainerWithLayout(
					layout.NewGridLayout(3),
					domain.Elapsed,
					domain.Data,
					domain.Files,
				),
			),
		)
	}

	scroll := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(1000, 400)))
	scroll.AddObject(widget.NewScrollContainer(completed))

	bottom := fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(1000, 40)),
		widget.NewButton("Close", func() {
			gui.app.Quit()
		}),
	)

	divider := fyne.NewContainerWithLayout(
		layout.NewGridLayout(1),
		widget.NewLabel(""),
		widget.NewLabelWithStyle("--- Completed Domains ---", fyne.TextAlignCenter, fyne.TextStyle{}),
	)

	content := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	content.AddObject(current)
	content.AddObject(done)
	content.AddObject(divider)
	content.AddObject(scroll)
	content.AddObject(bottom)
	content.Resize(fyne.NewSize(1000, 675))
	gui.window.SetContent(content)

	for _, domain := range gui.domains {
		domain.Name.Show()
		domain.Elapsed.Show()
		domain.ETA.Show()
		domain.Files.Show()
		domain.Data.Show()
		if !domain.Completed {
			domain.ProgressBar.Show()
		} else {
			domain.Files.Show()
			domain.ProgressBar.Hide()
		}
		domain.StatusLines[0].Show()
		domain.StatusLines[1].Show()
	}
}

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

// StartNewDomain will add a new DomainProgress to gui's domains, subsequent status
// and summary messages will refer to this domain until a new domain is started
func (gui *ProgressGUI) StartNewDomain(groupName string, domainName string) {
	domain := d.Load(
		domainName,
		groupName,
		fmt.Sprintf("%s/%s", paths.DFB(), groupName),
	)
	domainProgress := &DomainProgress{
		Domain:      domain,
		Name:        widget.NewLabel(fmt.Sprintf("backing up %s", domain.Name)),
		Elapsed:     widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{}),
		ETA:         widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{}),
		Files:       widget.NewLabelWithStyle("N/A", fyne.TextAlignTrailing, fyne.TextStyle{}),
		Data:        widget.NewLabelWithStyle("N/A", fyne.TextAlignTrailing, fyne.TextStyle{}),
		ProgressBar: widget.NewProgressBar(),
		StatusLines: []*widget.Label{
			widget.NewLabel(""),
			widget.NewLabel(""),
		},
	}
	domainProgress.ProgressBar.Max = 100
	gui.domains = append([]*DomainProgress{domainProgress}, gui.domains...)
	gui.currentDomain = domainProgress
	gui.updateLayout()
}

// StartNewEmptyDomain will add a new DomainProgress to gui's domains, however this
// domain won't be able to receive any status or summary messages
func (gui *ProgressGUI) StartNewEmptyDomain(groupName string, domainName string) {
	domainProgress := &DomainProgress{
		Domain:      d.Domain{},
		Completed:   false,
		Name:        widget.NewLabel(fmt.Sprintf("backing up %s", domainName)),
		Elapsed:     widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{}),
		ETA:         widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{}),
		Files:       widget.NewLabelWithStyle("N/A", fyne.TextAlignTrailing, fyne.TextStyle{}),
		Data:        widget.NewLabelWithStyle("N/A", fyne.TextAlignTrailing, fyne.TextStyle{}),
		ProgressBar: widget.NewProgressBar(),
		StatusLines: []*widget.Label{
			widget.NewLabel(""),
			widget.NewLabel(""),
		},
	}
	domainProgress.ProgressBar.Max = 100
	gui.domains = append(gui.domains, domainProgress)
	gui.currentDomain = domainProgress
	gui.updateLayout()
}
