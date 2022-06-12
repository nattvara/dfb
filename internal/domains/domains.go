package domains

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/nattvara/dfb/internal/paths"
)

// Domain is a directory a directory to backup
type Domain struct {
	Name          string   // Name of domain
	GroupName     string   // Name of group domain belongs to
	Path          string   // Path to domain eg. ~/domain ~/domain.somefile
	TemporaryPath string   // If Path does not exist a temporary path will be created, this might differ from Path
	ConfigPath    string   // Path to domain config ~/.dfb/[group]/domains/domain
	Repositories  string   // Comma separated list of repositories
	config        string   // Config content
	Symlink       *Symlink // Path to real domain src if it's a symlinked domain
}

// ParseConfig will parse the config stored at the domains Path
func (domain *Domain) ParseConfig() {
	file, err := ioutil.ReadFile(domain.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	domain.config = string(file)
	domain.parsePathFromConfig()
	domain.parseRepositoriesFromConfig()
	domain.parseSymlinkFromConfig()
}

func (domain *Domain) parsePathFromConfig() {
	var re = regexp.MustCompile(`(?m)path: (.*)\n`)
	domain.Path = re.FindStringSubmatch(domain.config)[1]
}

func (domain *Domain) parseRepositoriesFromConfig() {
	var re = regexp.MustCompile(`(?m)repos: (.*)\n`)
	domain.Repositories = re.FindStringSubmatch(domain.config)[1]
}

func (domain *Domain) parseSymlinkFromConfig() {
	var re = regexp.MustCompile(`(?m)symlink: (.*)\n`)

	matches := re.FindStringSubmatch(domain.config)
	if len(matches) < 2 {
		return
	}

	path := re.FindStringSubmatch(domain.config)[1]
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return
	}

	domain.Symlink = &Symlink{
		Source: path,
		Proxy:  fmt.Sprintf("%s/%s/symlinks/%s", paths.DFB(), domain.GroupName, domain.Name),
		Domain: domain,
	}
}

// Load domain will create domain type and load its config
func Load(name string, groupName string, groupPath string) Domain {
	domain := Domain{
		Name:       name,
		GroupName:  groupName,
		ConfigPath: fmt.Sprintf("%s/domains/%s", groupPath, name),
	}

	domain.ParseConfig()
	if domain.IsSingleFileDomain() {
		domain.TemporaryPath = domain.Path + ".dfb"
	}
	return domain
}

// SaveConfig saves the config for the Domain
func (domain *Domain) SaveConfig() {
	template := `path: %s
symlink: %s
exclusions: **/node_modules **/.DS_Store **/venv
repos: %s
`
	var symlink string
	if domain.Symlink != nil {
		symlink = domain.Symlink.Source
	}

	content := fmt.Sprintf(template, domain.Path, symlink, domain.Repositories)

	file, err := os.Create(domain.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_, err = file.WriteString(content)

	if err != nil {
		log.Fatal(err)
	}
}

// CreatePathIfNotCreated will create the writable path if not created
func (domain *Domain) CreatePathIfNotCreated() {
	if !domain.PathExists() && !domain.IsSingleFileDomain() {
		log.Printf("[domain: %s] creating temporary path", domain.Name)
		domain.CreatePath()
		domain.CreateTemporaryFlag()
	}
	if domain.IsSingleFileDomain() && !paths.Exists(domain.writablePath()) {
		log.Printf("[domain: %s] is single file domain, creating temporary path", domain.Name)
		domain.CreatePath()
		domain.CreateTemporaryFlag()
	}
}

// DeletePath will delete the domains writable path
func (domain *Domain) DeletePath() {
	os.RemoveAll(domain.writablePath())
}

func (domain *Domain) writablePath() string {
	if domain.TemporaryPath != "" {
		return domain.TemporaryPath
	}
	return domain.Path
}

// PathExists checks if the domains path exists
func (domain *Domain) PathExists() bool {
	return paths.Exists(domain.Path)
}

// WritablePathExists checks if the domains writable path exists
func (domain *Domain) WritablePathExists() bool {
	return paths.Exists(domain.writablePath())
}

// CreatePath creates the domains path
func (domain *Domain) CreatePath() {
	err := os.MkdirAll(domain.writablePath(), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// CreateTemporaryFlag writes a file to the domains path to flag it as temporary
func (domain *Domain) CreateTemporaryFlag() {
	err := ioutil.WriteFile(domain.TemporaryFlag(), []byte{}, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// DeleteTemporaryFlag will delete the the temporary flag in the domain
func (domain *Domain) DeleteTemporaryFlag() {
	err := os.Remove(domain.TemporaryFlag())
	if err != nil {
		log.Fatal(err)
	}
}

// TemporaryFlag returns the path to the flag that indicates that a directory is temporary
func (domain *Domain) TemporaryFlag() string {
	return fmt.Sprintf("%s/.dfb-temporary", domain.writablePath())
}

// IsSingleFileDomain checks if the domain is only for a single file
func (domain *Domain) IsSingleFileDomain() bool {
	if !paths.Exists(domain.Path) {
		// If path does not exist, treat as directory domain even if
		// backed up domain might be a single file domain
		return false
	}
	return !paths.IsDir(domain.Path)
}

// IsSymlinkedDomain check if the domain is a symlinked domain
// meaning that its content does not exist at path but at some
// other location
func (domain *Domain) IsSymlinkedDomain() bool {
	return domain.Symlink != nil
}

// IsTemporary checks whether domain is temporary
func (domain *Domain) IsTemporary() bool {
	return paths.Exists(domain.TemporaryFlag())
}

// IsEmpty checks if domain is empty of any files or directories
func (domain *Domain) IsEmpty() bool {
	files, err := filepath.Glob(domain.writablePath() + "/*")
	if err != nil {
		log.Fatal(err)
	}
	var count int
	for _, file := range files {
		basename := filepath.Base(file)
		if basename == ".DS_Store" {
			continue
		}
		if file == domain.TemporaryFlag() {
			continue
		}
		count++
	}
	return count == 0
}

// LinkToBackupsExist checks whether symlinks to the domains backups exist
func (domain *Domain) LinkToBackupsExist() bool {
	return paths.SymlinkExists(domain.backupLinkTargetPath())
}

// CreateLinkToBackups creates a symlink to domain backups in the mounted restic filesystem
func (domain *Domain) CreateLinkToBackups(mountpoint string) {
	os.Symlink(domain.backupLinkSourcePath(mountpoint), domain.backupLinkTargetPath())
}

// DeleteLinkToBackups deletes symlink to backups
func (domain *Domain) DeleteLinkToBackups(mountpoint string) {
	os.Remove(domain.backupLinkTargetPath())
}

func (domain *Domain) backupLinkTargetPath() string {
	if domain.IsSingleFileDomain() {
		return fmt.Sprintf("%s/__recover__", domain.writablePath())
	}
	return fmt.Sprintf("%s/__recover__", domain.Path)
}

func (domain *Domain) backupLinkSourcePath(mountpoint string) string {
	return fmt.Sprintf("%s/tags/%s", mountpoint, domain.Name)
}
