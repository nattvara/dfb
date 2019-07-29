package domains

import (
	"log"
	"os"

	"github.com/nattvara/dfb/internal/paths"
)

// Symlink is a symbolik link to domains real content, not all domains
// have symbolik links to their content
type Symlink struct {
	Source string  // The source of the domains symlink
	Proxy  string  // Domains with symlinks have a symlink in ~/.dfb/[group]/symlinks/[domain] to Source that is never deleted or repointed
	domain *Domain // Domain symlink belongs to
}

// Availible checks if the symlinks source content is availible
func (symlink *Symlink) Availible() bool {
	return paths.Exists(symlink.Source)
}

// LinkDomainToProxyIfNotLinked will link domain Path to symlink proxy
// if domain path does not exist and is not temporary
func (symlink *Symlink) LinkDomainToProxyIfNotLinked() {
	if !symlink.domain.PathExists() {
		log.Printf("[domain: %s] domain path did not exist, linking to source", symlink.domain.Name)
		os.Symlink(symlink.Proxy, symlink.domain.Path)
		return
	}
	if paths.Exists(symlink.domain.TemporaryFlag()) {
		log.Printf("[domain: %s] domain path was temporary, linking to source", symlink.domain.Name)
		os.RemoveAll(symlink.domain.Path)
		os.Symlink(symlink.Proxy, symlink.domain.Path)
	}
}

// UnlinkDomainFromProxyIfLinked will unlink domain from proxy will unlink
// domain from its proxy if linked
func (symlink *Symlink) UnlinkDomainFromProxyIfLinked() {
	if paths.IsSymlink(symlink.domain.Path) {
		log.Printf("[domain: %s] symlink source went away, removing symlink", symlink.domain.Name)
		os.Remove(symlink.domain.Path)
	}
}
