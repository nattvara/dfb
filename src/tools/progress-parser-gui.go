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

	NameWidget   *widget.Label
	ETA          *widget.Label
	ProgressBar  *widget.ProgressBar
	CurrentFiles []*canvas.Text
}

// LoadUI will load the initial UI for gui
func (gui *ProgressGUI) LoadUI(app fyne.App) {
	gui.window = app.NewWindow("dfb Progress Report")

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

	currentFiles := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	for _, file := range gui.currentDomain.CurrentFiles {
		file.Resize(fyne.NewSize(700, 10))
		currentFiles.AddObject(fyne.NewContainer(
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
		currentFiles,
		bottom,
	)
	content.Resize(fyne.NewSize(700, 475))
	gui.window.SetContent(content)

	for _, domain := range gui.domains {
		domain.NameWidget.Show()
		domain.ETA.Show()
		domain.ProgressBar.Show()
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
	gui.currentDomain.ETA.SetText(msg.GetETAString())
	gui.currentDomain.ProgressBar.SetValue(msg.GetProcent())

	if len(msg.CurrentFiles) == 1 {
		gui.currentDomain.CurrentFiles[0].Text = msg.CurrentFiles[0]
	} else if len(msg.CurrentFiles) == 2 {
		gui.currentDomain.CurrentFiles[0].Text = msg.CurrentFiles[0]
		gui.currentDomain.CurrentFiles[1].Text = msg.CurrentFiles[1]
	}
}

func (gui *ProgressGUI) handleSummaryMessage(msg restic.SummaryMessage) {
	gui.currentDomain.ETA.SetText(fmt.Sprintf(
		"took %s, added: %s",
		msg.GetDurationString(),
		msg.GetDataAddedString(),
	))
	gui.currentDomain.ProgressBar.SetValue(100)
	gui.currentDomain.CurrentFiles[0].Text = ""
	gui.currentDomain.CurrentFiles[1].Text = ""
}

func (gui *ProgressGUI) handleDFBMessage(msg restic.DFBMessage) {
	switch msg.Action {
	case "begin":
		gui.StartNewDomain(msg.Group, msg.Domain)
	case "unavailable":
		gui.StartNewEmptyDomain(msg.Group, msg.Domain)
		gui.currentDomain.ETA.SetText("unavailable")
		gui.currentDomain.ProgressBar.SetValue(100)
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
	file1 := canvas.NewText("", color.Gray{})
	file1.TextSize = 9
	file2 := canvas.NewText("", color.Gray{})
	file2.TextSize = 9
	domainProgress := &DomainProgress{
		Domain:      domain,
		NameWidget:  widget.NewLabel(fmt.Sprintf("backing up %s", domain.Name)),
		ETA:         widget.NewLabel("N/A"),
		ProgressBar: widget.NewProgressBar(),
		CurrentFiles: []*canvas.Text{
			file1,
			file2,
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
	DomainProgress := &DomainProgress{
		Domain:      d.Domain{},
		NameWidget:  widget.NewLabel(fmt.Sprintf("backing up %s", domainName)),
		ETA:         widget.NewLabel("N/A"),
		ProgressBar: widget.NewProgressBar(),
	}
	DomainProgress.ProgressBar.Max = 100
	gui.domains = append(gui.domains, DomainProgress)
	gui.currentDomain = DomainProgress
	gui.updateLayout()
}
