package progress

import (
	"fmt"

	d "github.com/nattvara/dfb/internal/domains"
	"github.com/nattvara/dfb/internal/paths"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

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
