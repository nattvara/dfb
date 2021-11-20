package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Bottom component is a container component for controls
type Bottom struct {
	container *fyne.Container

	App fyne.App
}

// GetNewContainer returns a container with controls
func (c *Bottom) GetNewContainer() *fyne.Container {
	c.container = fyne.NewContainerWithLayout(
		layout.NewGridLayout(1),
		widget.NewButton("Close", func() {
			c.App.Quit()
		}),
	)
	return c.container
}
