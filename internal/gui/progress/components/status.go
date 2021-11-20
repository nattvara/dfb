package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Status is a component with the status of the current domain being backed up
type Status struct {
	container *fyne.Container

	title       *widget.Label
	elapsed     *widget.Label
	eta         *widget.Label
	filesDone   *widget.Label
	bytesDone   *widget.Label
	progressBar *widget.ProgressBar
	statusLines []*widget.Label
}

// GetNewContainer returns a fyne container with the necessary widgets
// to display backup status
func (c *Status) GetNewContainer() *fyne.Container {
	c.title = widget.NewLabel("N/A")
	c.elapsed = widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{})
	c.eta = widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{})
	c.filesDone = widget.NewLabelWithStyle("N/A", fyne.TextAlignTrailing, fyne.TextStyle{})
	c.bytesDone = widget.NewLabelWithStyle("N/A", fyne.TextAlignTrailing, fyne.TextStyle{})
	c.progressBar = widget.NewProgressBar()
	c.progressBar.Max = 100
	c.statusLines = []*widget.Label{
		widget.NewLabel(""),
		widget.NewLabel(""),
	}

	c.container = fyne.NewContainerWithLayout(
		layout.NewGridLayout(1),
		fyne.NewContainerWithLayout(
			layout.NewGridLayout(1),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(1),
				c.title,
			),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(1),
				c.progressBar,
			),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(4),
				c.elapsed,
				c.eta,
				c.filesDone,
				c.bytesDone,
			),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(1),
				c.statusLines[0],
			),
			fyne.NewContainerWithLayout(
				layout.NewGridLayout(1),
				c.statusLines[1],
			),
		),
	)
	return c.container
}

// SetTitle sets the title of the status component
func (c *Status) SetTitle(title string) {
	c.title.SetText(title)
}

// SetElapsed sets the elapsed time for current domain snapshot
func (c *Status) SetElapsed(elapsed string) {
	c.elapsed.SetText("Elapsed: " + elapsed)
}

// SetETA sets the current ETA for domain snapshot
func (c *Status) SetETA(eta string) {
	c.eta.SetText("ETA: " + eta)
}

// SetFilesDone sets the number of files that have been processed for current snapshot
func (c *Status) SetFilesDone(done int, total int) {
	c.filesDone.SetText(fmt.Sprintf("Files: %v/%v", done, total))
}

// SetBytesDone sets the number of bytes that have been processed for current snapshot
func (c *Status) SetBytesDone(done string, total string) {
	c.bytesDone.SetText(fmt.Sprintf("%s/%s", done, total))
}

// SetProgress updates the status bar of component c
func (c *Status) SetProgress(procent float64) {
	c.progressBar.SetValue(procent)
}

// SetFirstStatusLine sets the first status line of component c
func (c *Status) SetFirstStatusLine(msg string) {
	c.statusLines[0].SetText(msg)
}

// SetSecondStatusLine sets the second status line of component c
func (c *Status) SetSecondStatusLine(msg string) {
	c.statusLines[1].SetText(msg)
}

// Clear clears all values of widgets in component c
func (c *Status) Clear() {
	c.title.SetText("")
	c.elapsed.SetText("")
	c.eta.SetText("")
	c.filesDone.SetText("")
	c.bytesDone.SetText("")
	c.progressBar.Hide()
	c.SetFirstStatusLine("")
	c.SetSecondStatusLine("")
}

// Reset resets the status to it's initial state
func (c *Status) Reset() {
	c.title.SetText("")
	c.elapsed.SetText("")
	c.eta.SetText("")
	c.filesDone.SetText("")
	c.bytesDone.SetText("")
	c.progressBar.SetValue(0)
	c.SetFirstStatusLine("")
	c.SetSecondStatusLine("")
}

// SetStatusLinesFromCurrentFiles updates status lines from slice of file paths
func (c *Status) SetStatusLinesFromCurrentFiles(currentFiles []string) {
	if len(currentFiles) == 1 {
		c.SetFirstStatusLine(truncateString(currentFiles[0], 100))
		c.SetSecondStatusLine("")
	} else if len(currentFiles) == 2 {
		c.SetFirstStatusLine(truncateString(currentFiles[0], 100))
		c.SetSecondStatusLine(truncateString(currentFiles[1], 100))
	} else {
		c.SetFirstStatusLine("")
		c.SetSecondStatusLine("")
	}
}

// truncateString truncates string str to max number of characters
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
