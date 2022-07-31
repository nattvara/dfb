package restic

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/nattvara/dfb/internal/paths"
)

// ExecuteKeyListCommand will execute the restic key list command for the repo at the given path,
// the command returns 0 if the password is correct, 1 if the password is invalid
// $ restic -r [REPO] key list
func ExecuteKeyListCommand(repoPath string, password string) int {
	app := "restic"
	arg0 := "-r"
	arg1 := repoPath
	arg2 := "key"
	arg3 := "list"
	arg4 := "--json"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4)

	buffer := bytes.Buffer{}
	buffer.Write([]byte(password))
	cmd.Stdin = &buffer

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
	}

	return 0
}

// ExecuteKeyListCommand will execute the restic backup command
// $ restic -r [REPO] backup [PATH] --tag [DOMAIN] --exclude-file [exclusions] --verbose --json
func ExecuteBackupCommand(
	repo string,
	password string,
	domainName string,
	workingDir string,
	backupPath string,
	exclusions string,
	msgCh chan string,
) (int, error) {
	exclusionsFilename := fmt.Sprintf("/tmp/dfb_exclusions_%s_%d", domainName, time.Now().Unix())
	if paths.Exists(exclusionsFilename) {
		os.Remove(exclusionsFilename)
	}

	exclusionsFile, err := os.Create(exclusionsFilename)
	if err != nil {
		log.Println(err)
		return 1, err
	}
	defer exclusionsFile.Close()

	_, err = exclusionsFile.WriteString(exclusions)
	if err != nil {
		log.Println(err)
		return 1, err
	}

	app := "restic"
	arg0 := "-r"
	arg1 := repo
	arg2 := "backup"
	arg3 := backupPath
	arg4 := "--tag"
	arg5 := domainName
	arg6 := "--exclude-file"
	arg7 := exclusionsFilename
	arg8 := "--verbose"
	arg9 := "--json"

	cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9)
	cmd.Dir = workingDir

	buffer := bytes.Buffer{}
	buffer.Write([]byte(password))
	cmd.Stdin = &buffer

	return executeCmdAndSendOutputToChannel(cmd, msgCh)
}

func executeCmdAndSendOutputToChannel(cmd *exec.Cmd, msgCh chan string) (int, error) {
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout // Use the same pipe for stderr

	wg := new(sync.WaitGroup)
	wg.Add(1)

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			msgCh <- line
		}

		wg.Done()
	}()

	err := cmd.Start()
	if err != nil {
		log.Println(err)
		return 1, err
	}

	// Ensures that both the cmd AND the output scanner has finished
	wg.Wait()
	err = cmd.Wait()

	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode(), nil
	}

	return 0, nil
}
