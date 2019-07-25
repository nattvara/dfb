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
	d "dfb/src/internal/domains"
	"dfb/src/internal/paths"
	"dfb/src/internal/restic"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
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

	window fyne.Window
	app    fyne.App
}

// DomainProgress contains a domain and the widgets to display the backup progress
// of that domain
type DomainProgress struct {
	Domain d.Domain

	NameWidget  *widget.Label
	ETA         *widget.Label
	ProgressBar *widget.ProgressBar
	StatusLines []*canvas.Text
}

// LoadUI will load the initial UI for gui
func (gui *ProgressGUI) LoadUI(app fyne.App) {
	now := time.Now().Format("15:04")
	gui.window = app.NewWindow("Progress report for dfb backup started at " + now)

	gui.window.SetContent(widget.NewLabel("waiting for messages on stdin"))
	gui.window.Show()
}

func (gui *ProgressGUI) updateLayout() {
	domains := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	for _, domain := range gui.domains {
		domains.AddObject(
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(3),
				domain.NameWidget,
				domain.ETA,
				domain.ProgressBar,
			),
		)
	}

	scroll := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(700, 400)))
	scroll.AddObject(widget.NewScrollContainer(domains))

	statusLines := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	for _, file := range gui.currentDomain.StatusLines {
		file.Resize(fyne.NewSize(700, 10))
		statusLines.AddObject(fyne.NewContainer(
			file,
		))
	}

	bottom := fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(700, 40)),
		widget.NewButton("Close", func() {
			gui.app.Quit()
		}),
	)

	content := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		scroll,
		statusLines,
		bottom,
	)
	content.Resize(fyne.NewSize(700, 475))
	gui.window.SetContent(content)

	for _, domain := range gui.domains {
		domain.NameWidget.Show()
		domain.ETA.Show()
		domain.ProgressBar.Show()
	}
	gui.currentDomain.StatusLines[0].Show()
	gui.currentDomain.StatusLines[1].Show()
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
	gui.currentDomain.ETA.SetText(msg.GetETAString())
	gui.currentDomain.ProgressBar.SetValue(msg.GetProcent())

	if len(msg.CurrentFiles) == 1 {
		gui.currentDomain.StatusLines[0].Text = msg.CurrentFiles[0]
	} else if len(msg.CurrentFiles) == 2 {
		gui.currentDomain.StatusLines[0].Text = msg.CurrentFiles[0]
		gui.currentDomain.StatusLines[1].Text = msg.CurrentFiles[1]
	}
}

func (gui *ProgressGUI) handleSummaryMessage(msg restic.SummaryMessage) {
	gui.currentDomain.ETA.SetText(fmt.Sprintf(
		"took %s, added: %s",
		msg.GetDurationString(),
		msg.GetDataAddedString(),
	))
	gui.currentDomain.ProgressBar.SetValue(100)
	gui.currentDomain.StatusLines[0].Text = ""
	gui.currentDomain.StatusLines[1].Text = ""
}

func (gui *ProgressGUI) handleDFBMessage(msg restic.DFBMessage) {
	switch msg.Action {
	case "begin":
		gui.StartNewDomain(msg.Group, msg.Domain)
	case "unavailable":
		gui.StartNewEmptyDomain(msg.Group, msg.Domain)
		gui.currentDomain.ETA.SetText("unavailable")
		gui.currentDomain.ProgressBar.SetValue(100)
	case "gathering_stats":
		if gui.currentDomain == nil {
			return
		}
		gui.currentDomain.StatusLines[0].Text = "gathering stats for " + msg.Domain
		gui.updateLayout()
	case "gathering_stats_done":
		if gui.currentDomain == nil {
			return
		}
		gui.currentDomain.StatusLines[1].Text = "done."
		gui.updateLayout()
	default:
		gui.StartNewEmptyDomain(msg.Group, msg.Domain)
		gui.currentDomain.ETA.SetText(msg.Action)
		gui.currentDomain.ProgressBar.SetValue(100)
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
	line1 := canvas.NewText("", color.White)
	line1.TextSize = 9
	line2 := canvas.NewText("", color.White)
	line2.TextSize = 9
	domainProgress := &DomainProgress{
		Domain:      domain,
		NameWidget:  widget.NewLabel(fmt.Sprintf("backing up %s", domain.Name)),
		ETA:         widget.NewLabel("N/A"),
		ProgressBar: widget.NewProgressBar(),
		StatusLines: []*canvas.Text{
			line1,
			line2,
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
	line1 := canvas.NewText("", color.White)
	line1.TextSize = 9
	line2 := canvas.NewText("", color.White)
	line2.TextSize = 9
	domainProgress := &DomainProgress{
		Domain:      d.Domain{},
		NameWidget:  widget.NewLabel(fmt.Sprintf("backing up %s", domainName)),
		ETA:         widget.NewLabel("N/A"),
		ProgressBar: widget.NewProgressBar(),
		StatusLines: []*canvas.Text{
			line1,
			line2,
		},
	}
	domainProgress.ProgressBar.Max = 100
	gui.domains = append(gui.domains, domainProgress)
	gui.currentDomain = domainProgress
	gui.updateLayout()
}
