package progress

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

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
