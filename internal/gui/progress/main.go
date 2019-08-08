package progress

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

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
