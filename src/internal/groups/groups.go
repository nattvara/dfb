package groups

import (
	"fmt"
	"io/ioutil"
	"log"

	d "dfb/src/internal/domains"
	"dfb/src/internal/paths"
)

// FetchGroups reads and returns the groups stored on disk in dfp path
func FetchGroups() []Group {
	files, err := ioutil.ReadDir(paths.DFB())
	if err != nil {
		log.Fatal(err)
	}

	var groups []Group

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		groups = append(groups, Group{
			Name: f.Name(),
			Path: fmt.Sprintf("%s/%s", paths.DFB(), f.Name()),
		})
	}

	return groups
}

// NumberOfGroupsMounted counts number of mounted groups
func NumberOfGroupsMounted(groups []Group) int {
	count := 0
	for _, group := range groups {
		if group.IsMounted() {
			count++
		}
	}
	return count
}

// Group contians a number of domains, and restic repositories
// to backup those domains to.
type Group struct {
	Path string
	Name string
}

// Mountpoint returns the path to where the group will mount restic repos
func (group *Group) Mountpoint() string {
	return fmt.Sprintf("%s/mountpoint", group.Path)
}

// IsMounted checks if the group have been mounted
func (group *Group) IsMounted() bool {
	dir, err := ioutil.ReadDir(group.Mountpoint())
	if err != nil {
		return false
	}

	if len(dir) > 0 {
		return true
	}
	return false
}

// Domains returns the domains belonging to the group
func (group *Group) Domains() []d.Domain {
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/domains", group.Path))
	if err != nil {
		log.Fatal(err)
	}

	var domains []d.Domain
	for _, file := range files {
		domain := d.Load(file.Name(), group.Name, group.Path)
		domains = append(domains, domain)
	}

	return domains
}
