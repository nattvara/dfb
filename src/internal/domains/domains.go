package domains

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"dfb/src/internal/paths"
)

// Domain is a directory a directory to backup
type Domain struct {
	Name          string // Name of domain
	Path          string // Path to domain eg. ~/domain ~/domain.somefile
	TemporaryPath string // If Path does not exist a temporary path will be created, this might differ from Path
	ConfigPath    string // Path to domain config ~/.dfb/[group]/domains/domain
	config        string // Config content
	Symlink       string // Path to real domain src if it's a symlinked domain
}

// ParseConfig will parse the config stored at the domains Path
func (domain *Domain) ParseConfig() {
	file, err := ioutil.ReadFile(domain.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	domain.config = string(file)
	domain.parsePathFromConfig()
	domain.parseSymlinkFromConfig()
}

func (domain *Domain) parsePathFromConfig() {
	var re = regexp.MustCompile(`(?m)path: (.*)\n`)
	domain.Path = re.FindStringSubmatch(domain.config)[1]
}

func (domain *Domain) parseSymlinkFromConfig() {
	var re = regexp.MustCompile(`(?m)symlink: (.*)\n`)

	matches := re.FindStringSubmatch(domain.config)
	if len(matches) < 2 {
		return
	}

	domain.Symlink = re.FindStringSubmatch(domain.config)[1]
}

// Load domain will create domain type and load its config
func Load(name string, groupPath string) Domain {
	domain := Domain{
		Name:       name,
		ConfigPath: fmt.Sprintf("%s/domains/%s", groupPath, name),
	}

	domain.ParseConfig()
	if domain.IsSingleFileDomain() {
		domain.TemporaryPath = domain.Path + ".dfb"
	}
	return domain
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

// DeletePathIfTemporary will delete the domains writable path if it's temporary
func (domain *Domain) DeletePathIfTemporary() {
	if paths.Exists(domain.temporaryFlag()) {
		log.Printf("[domain: %s] removing temporary path", domain.Name)
		os.RemoveAll(domain.writablePath())
	}
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

// CreatePath creates the domains path
func (domain *Domain) CreatePath() {
	err := os.MkdirAll(domain.writablePath(), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// CreateTemporaryFlag writes a file to the domains path to flag it as temporary
func (domain *Domain) CreateTemporaryFlag() {
	err := ioutil.WriteFile(domain.temporaryFlag(), []byte{}, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (domain *Domain) temporaryFlag() string {
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
