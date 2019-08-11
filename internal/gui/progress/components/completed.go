package components

import (
	"fmt"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

const (
	// ListLength is how many items should be showing at any time
	ListLength = 10
)

// Completed is a component with a paginated list of completed domain snapshots
type Completed struct {
	container *fyne.Container

	list         *fyne.Container
	items        []*item
	domains      []map[string]string
	addedDomains map[string]int

	controls    *fyne.Container
	next        *widget.Button
	previous    *widget.Button
	label       *widget.Label
	currentPage int
	pageCount   int
}

// GetNewContainer returns a fyne container with the necessary widgets
// to display list of snapshots and various metrics
func (c *Completed) GetNewContainer() *fyne.Container {
	c.addedDomains = make(map[string]int)
	c.next = widget.NewButton("Next", c.Next)
	c.previous = widget.NewButton("Previous", c.Previous)
	c.label = widget.NewLabelWithStyle(
		"1/1",
		fyne.TextAlignCenter,
		fyne.TextStyle{},
	)
	c.controls = fyne.NewContainerWithLayout(
		layout.NewGridLayout(3),
		widget.NewLabel(""),
		fyne.NewContainerWithLayout(
			layout.NewGridLayout(3),
			c.previous,
			c.label,
			c.next,
		),
	)
	c.list = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
	)
	c.container = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		c.controls,
		c.list,
	)
	for i := 1; i <= ListLength; i++ {
		item := &item{}
		c.items = append(c.items, item)
		c.list.AddObject(item.GetContainer())
	}
	return c.container
}

// AddCompletedDomain adds a completed domain to the list of completed domains,
// if necessary the list will increment to not show more items than ListLength
func (c *Completed) AddCompletedDomain(name string, duration string, processed string, added string) {
	if _, ok := c.addedDomains[name]; ok {
		return
	}
	c.addedDomains[name] = len(c.domains)

	domain := map[string]string{
		"name":      fmt.Sprintf("%v %s", len(c.domains)+1, name),
		"duration":  duration,
		"processed": processed,
		"added":     added,
	}
	c.domains = append(c.domains, domain)

	shouldUpdateList := true

	// Only store the domain if user have moved page in the list,
	// app is tricky to use otherwise
	if c.currentPage != c.pageCount {
		shouldUpdateList = false
	}

	listIndex := int(math.Mod(float64(len(c.domains)-1), float64(ListLength)))
	if listIndex == 0 && len(c.domains) > ListLength-1 {
		c.pageCount++
		if shouldUpdateList {
			c.currentPage++
			c.clearList()
		}
		c.updateLabel()
	}

	if !shouldUpdateList {
		return
	}

	c.items[listIndex].SetName(domain["name"])
	c.items[listIndex].SetDuration(domain["duration"])
	c.items[listIndex].SetProcessed(domain["processed"])
	c.items[listIndex].SetAdded(domain["added"])
}

// clearList clears the items in the list of any content
func (c *Completed) clearList() {
	for i := 0; i < ListLength; i++ {
		c.items[i].Clear()
	}
}

// updateLabel updates the label for the list that shows current page and
// total number of pages
func (c *Completed) updateLabel() {
	c.label.SetText(fmt.Sprintf("%v/%v", c.currentPage+1, c.pageCount+1))
}

// Previous will "go back" in the list of domain snapshots
func (c *Completed) Previous() {
	if c.currentPage <= 0 {
		return
	}
	c.currentPage--
	c.clearList()
	start := (c.currentPage+1)*ListLength - 1
	for i := start; i > start-ListLength; i-- {
		if i > len(c.domains)-1 {
			continue
		}

		listIndex := int(math.Mod(float64(i), float64(ListLength)))
		c.items[listIndex].SetName(c.domains[i]["name"])
		c.items[listIndex].SetDuration(c.domains[i]["duration"])
		c.items[listIndex].SetProcessed(c.domains[i]["processed"])
		c.items[listIndex].SetAdded(c.domains[i]["added"])
	}
	c.updateLabel()
}

// Next will "go forward" in the list of domain snapshots
func (c *Completed) Next() {
	if c.currentPage >= c.pageCount {
		return
	}
	c.currentPage++
	c.clearList()
	start := (c.currentPage+1)*ListLength - 1
	for i := start; i > start-ListLength; i-- {
		if i > len(c.domains)-1 {
			continue
		}
		listIndex := int(math.Mod(float64(i), float64(ListLength)))
		c.items[listIndex].SetName(c.domains[i]["name"])
		c.items[listIndex].SetDuration(c.domains[i]["duration"])
		c.items[listIndex].SetProcessed(c.domains[i]["processed"])
		c.items[listIndex].SetAdded(c.domains[i]["added"])
	}
	c.updateLabel()
}

type item struct {
	name      *widget.Label
	duration  *widget.Label
	processed *widget.Label
	added     *widget.Label
}

// GetContainer returns a new container for item i
func (i *item) GetContainer() *fyne.Container {
	i.name = widget.NewLabel("")
	i.duration = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
	i.processed = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
	i.added = widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})
	return fyne.NewContainerWithLayout(
		layout.NewGridLayout(2),
		i.name,
		fyne.NewContainerWithLayout(
			layout.NewGridLayout(3),
			i.duration,
			i.processed,
			i.added,
		),
	)
}

// SetName sets the name of item i
func (i *item) SetName(name string) {
	i.name.SetText(name)
}

// SetDuration sets the duration of item i
func (i *item) SetDuration(duration string) {
	i.duration.SetText(duration)
}

// SetProcessed sets the bytes processed for item i
func (i *item) SetProcessed(processed string) {
	i.processed.SetText("Processed: " + processed)
}

// SetAdded sets the bytes added for item i
func (i *item) SetAdded(added string) {
	i.added.SetText("Added: " + added)
}

// Clear clears the list item from any values
func (i *item) Clear() {
	i.name.SetText("")
	i.duration.SetText("")
	i.processed.SetText("")
	i.added.SetText("")
}
