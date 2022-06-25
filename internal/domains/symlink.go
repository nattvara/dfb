package domains

import (
	"fmt"
	"log"
	"os"

	"github.com/nattvara/dfb/internal/paths"
)

// Symlink is a symbolik link to domains real content, not all domains
// have symbolik links to their content
type Symlink struct {
	Source string  // The source of the domains symlink
	Proxy  string  // Domains with symlinks have a symlink in ~/.dfb/[group]/symlinks/[domain] to Source that is never deleted or repointed
	Domain *Domain // Domain symlink belongs to
}

// CreateProxy creates the proxy for the symlink
func (symlink *Symlink) CreateProxy() error {
	if !symlink.Availible() {
		return fmt.Errorf("could not find source for proxy at: %s", symlink.Source)
	}
	os.Symlink(symlink.Source, symlink.Proxy)
	return nil
}

// DeleteProxy deletes the symlink proxy
func (symlink *Symlink) DeleteProxy() error {
	return os.Remove(symlink.Proxy)
}

// Availible checks if the symlinks source content is availible
func (symlink *Symlink) Availible() bool {
	return paths.Exists(symlink.Source)
}

// LinkDomainToProxyIfNotLinked will link domain Path to symlink proxy
// if domain path does not exist and is not temporary
func (symlink *Symlink) LinkDomainToProxyIfNotLinked() {
	if !symlink.Domain.PathExists() {
		log.Printf("[domain: %s] domain path did not exist, linking to source", symlink.Domain.Name)
		os.Symlink(symlink.Proxy, symlink.Domain.Path)
		return
	}
	if paths.Exists(symlink.Domain.TemporaryFlag()) {
		log.Printf("[domain: %s] domain path was temporary, linking to source", symlink.Domain.Name)
		os.RemoveAll(symlink.Domain.Path)
		os.Symlink(symlink.Proxy, symlink.Domain.Path)
	}
}

// UnlinkDomainFromProxyIfLinked will unlink domain from proxy will unlink
// domain from its proxy if linked
func (symlink *Symlink) UnlinkDomainFromProxyIfLinked() {
	if paths.IsSymlink(symlink.Domain.Path) {
		log.Printf("[domain: %s] symlink source went away, removing symlink", symlink.Domain.Name)
		os.Remove(symlink.Domain.Path)
	}
}
