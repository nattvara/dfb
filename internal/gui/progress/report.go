package progress

import (
	"fmt"
	"time"

	d "github.com/nattvara/dfb/internal/domains"
	"github.com/nattvara/dfb/internal/gui/progress/components"
	"github.com/nattvara/dfb/internal/paths"
	"github.com/nattvara/dfb/internal/restic"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
)

// New creates and returns a new progress report
func New(app fyne.App) *Report {
	report := &Report{}
	report.app = app

	whitelist := map[string]struct{}{
		"":                         {},
		"read password from stdin": {},
	}
	report.MessageReceiver = &MessageReceiver{
		Report:            report,
		WhitelistedErrors: whitelist,
	}
	return report
}

// Report is an OpenGL based gui app bwith the backup status of domains
type Report struct {
	domains         []*d.Domain
	MessageReceiver *MessageReceiver

	StatusComponent    *components.Status
	CompletedComponent *components.Completed
	BottomComponent    *components.Bottom

	window fyne.Window
	app    fyne.App
}

// LoadUI will load the initial UI for report
func (report *Report) LoadUI(app fyne.App) {
	now := time.Now().Format("15:04")
	report.window = app.NewWindow("Progress report for dfb backup started at " + now)

	report.StatusComponent = &components.Status{}
	report.CompletedComponent = &components.Completed{}
	report.BottomComponent = &components.Bottom{App: report.app}

	content := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	content.AddObject(report.StatusComponent.GetNewContainer())
	content.AddObject(report.CompletedComponent.GetNewContainer())
	content.AddObject(report.BottomComponent.GetNewContainer())
	report.window.SetContent(content)

	report.StatusComponent.SetTitle("Waiting for json messages on stdin")

	report.window.Show()
}

// StartNewDomain starts the backup of given domain of given group
func (report *Report) StartNewDomain(groupName string, domainName string) {
	report.StatusComponent.Reset()
	domain := d.Load(
		domainName,
		groupName,
		fmt.Sprintf("%s/%s", paths.DFB(), groupName),
	)
	report.domains = append(report.domains, &domain)
	report.StatusComponent.SetTitle(domain.Name)
}

// CompleteCurrentDomain completes the backup of the last started domain,
// metadata is provided with given SummaryMessage msg
func (report *Report) CompleteCurrentDomain(msg restic.SummaryMessage) {
	domain := report.domains[len(report.domains)-1]
	report.CompletedComponent.AddCompletedDomain(
		domain.Name,
		"Took: "+msg.GetDurationString(),
		msg.GetDataProcessedString(),
		msg.GetDataAddedString(),
	)
}

// CompleteUnavailibleDomain completes the snapshot of a domain that was unavailable
func (report *Report) CompleteUnavailibleDomain(group string, domain string, msg string) {
	report.CompletedComponent.AddCompletedDomain(
		domain,
		msg,
		"0 B",
		"0 B",
	)
}

// Done displays a message that the backup is done
func (report *Report) Done() {
	report.StatusComponent.Clear()
	report.StatusComponent.SetTitle("Backup is done.")
}
