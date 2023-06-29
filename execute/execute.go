package execute

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
)

const (
	timeConstant = 100
)

type Executor interface {
	// Errors writes command errors to stderr during execution.
	Errors(cmd *exec.Cmd, targetDir string, action string) error
	// Output returns the output of the command as a string.
	Output(cmd *exec.Cmd, targetDir string, action string) (string, error)
	// CommandExists returns true if the command exists in the path.
	CommandExists(cmd string) bool
}

func NewOsExecutor() Executor {
	return OsExecutor{}
}

// OsExecutor implements Executor using the os/exec package.
type OsExecutor struct{}

// Action string is used to log command info and wrap any returned errors
func (OsExecutor) Errors(cmd *exec.Cmd, targetDir string, action string) error {
	log.Println(action)
	cmd.Dir = targetDir
	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("%s: piping standard error for %q: %w", action, cmd.String(), err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%s: executing command %q: %w", action, cmd.String(), err)
	}

	s := spinner.New(spinner.CharSets[9], timeConstant*time.Millisecond)
	s.Prefix = fmt.Sprintf("%s: Waiting for command %q ", action, cmd.String())
	s.Start()

	if _, err := io.Copy(os.Stderr, cmdErr); err != nil {
		fmt.Fprintf(os.Stderr, "error copying command stderr\n")
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s: command %q finished with error: %w", action, cmd.String(), err)
	}

	s.Stop()
	return nil
}

func (OsExecutor) Output(cmd *exec.Cmd, targetDir string, action string) (string, error) {
	cmd.Dir = targetDir

	s := spinner.New(spinner.CharSets[9], timeConstant*time.Millisecond)
	s.Prefix = fmt.Sprintf("%s: Waiting for command %q ", action, cmd.String())
	s.Start()

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s: executing command %q: %w", action, cmd.String(), err)
	}
	s.Stop()
	return string(output), nil
}

func (OsExecutor) CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
