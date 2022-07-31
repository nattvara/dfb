package groups

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	d "github.com/nattvara/dfb/internal/domains"
	"github.com/nattvara/dfb/internal/paths"
)

// GetGroupFromString checks that the provided string is a valid group and returns that Group
func GetGroupFromString(name string) (*Group, error) {
	groups, err := FetchGroups()
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.Name == name {
			return &group, nil
		}
	}
	return nil, fmt.Errorf("could not find group %s", name)
}

// FetchGroups reads and returns the groups stored on disk in dfp path
func FetchGroups() ([]Group, error) {
	files, err := ioutil.ReadDir(paths.DFB())
	if err != nil {
		return nil, err
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

	return groups, nil
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

// New returns a new Group with given name
func New(name string) *Group {
	return &Group{
		Name: name,
		Path: fmt.Sprintf("%s/%s", paths.DFB(), name),
	}
}

// Create creates the necessary directories for a group
func (group *Group) Create() {
	os.Mkdir(group.Path, 0760)
	os.Mkdir(group.Path+"/repos", 0760)
	os.Mkdir(group.Path+"/domains", 0760)
	os.Mkdir(group.Path+"/symlinks", 0760)
	os.Mkdir(group.Path+"/stats", 0760)
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
func (group *Group) Domains() ([]d.Domain, error) {
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/domains", group.Path))
	if err != nil {
		return nil, err
	}

	var domains []d.Domain
	for _, file := range files {
		domain := d.Load(file.Name(), group.Name, group.Path)
		domains = append(domains, domain)
	}

	return domains, nil
}

// DomainsMap returns the domains belonging to the group as a map
func (group *Group) DomainsMap() (map[string]d.Domain, error) {
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/domains", group.Path))
	if err != nil {
		return nil, err
	}

	domains := make(map[string]d.Domain)
	for _, file := range files {
		domain := d.Load(file.Name(), group.Name, group.Path)
		domains[file.Name()] = domain
	}

	return domains, nil
}

// DomainExists returns boolean if domain exists inside the given Group
func (group *Group) DomainExists(domain string) bool {
	domains, err := group.DomainsMap()
	if err != nil {
		log.Fatalf(err.Error())
	}

	if _, ok := domains[domain]; ok {
		return true
	}
	return false
}

// AddDomainWithName adds a Domain with domainName
func (group *Group) AddDomainWithName(domainName string) error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	domain := d.Domain{
		Name:         domainName,
		GroupName:    group.Name,
		Path:         fmt.Sprintf("%s/%s", homedir, domainName),
		ConfigPath:   fmt.Sprintf("%s/domains/%s", group.Path, domainName),
		Repositories: AllRepositories,
	}

	domain.SaveConfig()

	return nil
}

// AddDomainWithNameAndSymlink adds a Domain with domainName and a Symlink
func (group *Group) AddDomainWithNameAndSymlink(domainName string, symlink string) error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	s := &d.Symlink{
		Source: symlink,
		Proxy:  fmt.Sprintf("%s/%s/symlinks/%s", paths.DFB(), group.Name, domainName),
	}

	domain := d.Domain{
		Name:         domainName,
		GroupName:    group.Name,
		Path:         fmt.Sprintf("%s/%s", homedir, domainName),
		ConfigPath:   fmt.Sprintf("%s/domains/%s", group.Path, domainName),
		Repositories: AllRepositories,
		Symlink:      s,
	}

	s.Domain = &domain

	s.CreateProxy()
	domain.SaveConfig()

	return nil
}

// Repositories reads the configured repositories of the given group
func (group *Group) Repositories() ([]*Repository, error) {
	var repos []*Repository

	path := group.Path + "/repos"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		content, err := ioutil.ReadFile(path + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}

		repos = append(repos, &Repository{
			Name:       f.Name(),
			ResticPath: string(content),
		})
	}

	return repos, nil
}

// GetRepositoryByName returns the Repository from Group g with name if it exists
func (g *Group) GetRepositoryByName(name string) (*Repository, error) {
	repos, err := g.Repositories()
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		if repo.Name == name {
			return repo, nil
		}
	}

	return nil, fmt.Errorf("did not find a repo with the name %s", name)
}

// AddRepository adds a Repository r to the config of Group g
func (g *Group) AddRepository(r Repository) error {
	repoConfigPath := fmt.Sprintf("%s/repos/%s", g.Path, r.Name)

	file, err := os.Create(repoConfigPath)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(r.ResticPath)

	if err != nil {
		return err
	}

	return nil
}
