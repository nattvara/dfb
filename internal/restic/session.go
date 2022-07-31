package restic

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/nattvara/dfb/internal/domains"
	"golang.org/x/crypto/ssh/terminal"
)

// Session is a type that represents a session against a single restic repository
type Session struct {
	RepositoryPath string
	RepositoryName string
	password       string
}

// NewSession creates a new session for repo with name and path
func NewSession(name string, path string) *Session {
	return &Session{
		RepositoryName: name,
		RepositoryPath: path,
	}
}

// PromptForPassword will ask the user to enter the repository password
func (s *Session) PromptForPassword() {
	fmt.Printf("Password: ")
	password, _ := terminal.ReadPassword(0)
	s.password = string(password)
	fmt.Print("\n")
}

// CheckPassword ensures password is correct
func (s *Session) CheckPassword() bool {
	res := ExecuteKeyListCommand(s.RepositoryPath, s.password)

	// exit code was 0, password was valid, otherwise invalid
	return res == 0
}

// StartDomainBackupConsumer starts a consumer for Domain types on the in channel to make a backup of
func (s *Session) StartDomainBackupConsumer(id int, in <-chan domains.Domain, out chan domains.Domain, msgCh chan string) {
	for domain := range in {
		var workingDirPath string
		var backupPath string
		if domain.IsSymlinkedDomain() {
			workingDirPath = domain.Symlink.Proxy
		} else {
			workingDirPath = domain.Path
		}

		if _, err := os.Stat(workingDirPath); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("%s domain unavailable\n", domain.Name)
			out <- domain
			continue
		}

		isDir, err := isDirectory(workingDirPath)
		if err != nil {
			log.Println(err)
			out <- domain
			continue
		}

		if !isDir {
			workingDirPath, backupPath = filepath.Split(workingDirPath)
		} else {
			backupPath = "."
		}

		ExecuteBackupCommand(
			s.RepositoryPath,
			s.password,
			domain.Name,
			workingDirPath,
			backupPath,
			domain.Exclusions,
			msgCh,
		)

		out <- domain
	}
}

// StartDomainStatisticsConsumer starts a consumer for Domain types on the in channel to check the stats for
func (s *Session) StartDomainStatisticsConsumer(id int, ch <-chan domains.Domain, msgCh chan string, wg *sync.WaitGroup) {
	for _ = range ch {
		wg.Done()
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}
