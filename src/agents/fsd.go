package main

import (
	g "dfb/src/internal/groups"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		os.Exit(0)
	}()

	groups := g.FetchGroups()
	groupsMounted := g.NumberOfGroupsMounted(groups)

	for _, group := range groups {

		domains := group.Domains()
		for _, domain := range domains {
			if !domain.IsSymlinkedDomain() {
				continue
			}
			if domain.Symlink.Availible() {
				domain.Symlink.LinkDomainToProxyIfNotLinked()
			} else {
				domain.Symlink.UnlinkDomainFromProxyIfLinked()
			}
		}

		if group.IsMounted() {
			for _, domain := range domains {
				domain.CreatePathIfNotCreated()
				if !domain.LinkToBackupsExist() {
					log.Printf("[domain: %s] has no link to backups, creating", domain.Name)
					domain.CreateLinkToBackups(group.Mountpoint())
				}
			}
		} else if groupsMounted == 0 {
			for _, domain := range domains {
				if domain.LinkToBackupsExist() {
					log.Printf("[domain: %s] has link to backups, removing", domain.Name)
					domain.DeleteLinkToBackups(group.Mountpoint())
				}
				domain.DeletePathIfTemporary()
			}
		}
	}
}
