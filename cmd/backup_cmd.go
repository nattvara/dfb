package cmd

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/nattvara/dfb/internal/domains"
	"github.com/nattvara/dfb/internal/groups"
	"github.com/nattvara/dfb/internal/output"
	"github.com/nattvara/dfb/internal/restic"
	"github.com/spf13/cobra"
)

const DomainStatisticsConsumers = 2
const DomainBackupConsumers = 2

var backupCmd = &cobra.Command{
	Use:   "backup [group] [repo]",
	Short: "Backup a group of domains to a repo",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		group, err := groups.GetGroupFromString(args[0])
		if err != nil {
			log.Println(err.Error())
			return errors.New("failed to list repos")
		}

		repo, err := group.GetRepositoryByName(args[1])
		if err != nil {
			log.Println(err.Error())
			return fmt.Errorf("failed to get repo with name %s, this most likely mean it is not configured for the given group", args[1])
		}

		r := restic.NewSession(repo.Name, repo.ResticPath)
		r.PromptForPassword()
		valid := r.CheckPassword()

		if !valid {
			return fmt.Errorf("invalid password provided for repository %s", repo.Name)
		}

		fmt.Println("password was valid")

		dmns, err := group.Domains()
		if err != nil {
			log.Println(err.Error())
			return fmt.Errorf("failed to read domains")
		}

		dmns = filterDomainsByRepo(dmns, repo)

		backupCh := make(chan domains.Domain)
		statsCh := make(chan domains.Domain)
		msgCh := make(chan string)

		// Producer and consumer wait groups
		wgp := new(sync.WaitGroup)
		wgc := new(sync.WaitGroup)
		wgp.Add(len(dmns))
		wgc.Add(len(dmns))

		// Consumer functions
		for i := 0; i < DomainBackupConsumers; i++ {
			go r.StartDomainBackupConsumer(i, backupCh, statsCh, msgCh)
		}
		for i := 0; i < DomainStatisticsConsumers; i++ {
			go r.StartDomainStatisticsConsumer(i, statsCh, msgCh, wgc)
		}

		// Producer function
		go func() {
			for _, domain := range dmns {
				backupCh <- domain
				wgp.Done()
			}
		}()

		// Output consumer
		go output.StartJsonMessageParser(msgCh)

		// Wait for all domains to be finished.
		wgp.Wait()
		wgc.Wait()

		close(backupCh)
		close(statsCh)

		return nil
	},
}

func filterDomainsByRepo(dmns []domains.Domain, repo *groups.Repository) []domains.Domain {
	var out []domains.Domain
	for _, domain := range dmns {
		if domain.Repositories == groups.AllRepositories {
			out = append(out, domain)
		} else if domain.RepositoriesContain(repo.Name) {
			out = append(out, domain)
		}
	}
	return out
}

func init() {
	RootCmd.AddCommand(backupCmd)
}
