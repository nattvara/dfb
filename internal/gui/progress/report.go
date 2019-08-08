package progress

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// New creates a progress report report
func New(app fyne.App) *Report {
	report := &Report{}
	report.app = app
	return report
}

// Report contains fyne and app related values
type Report struct {
	domains       []*DomainProgress
	currentDomain *DomainProgress
	done          bool

	window fyne.Window
	app    fyne.App
}

// LoadUI will load the initial UI for report
func (report *Report) LoadUI(app fyne.App) {
	now := time.Now().Format("15:04")
	report.window = app.NewWindow("Progress report for dfb backup started at " + now)

	report.window.SetContent(widget.NewLabel("waiting for messages on stdin"))
	report.window.Show()
}
