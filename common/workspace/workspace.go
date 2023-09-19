package workspace

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/vision-cli/vision/common/execute"
	"github.com/vision-cli/vision/common/file"
)

// Remove removes any go.work or go.work.sum files in projectDir.
func Remove(projectDir string) error {
	if err := file.RemoveNamed(projectDir, "go.work", "go.work.sum"); err != nil {
		return fmt.Errorf("removing existing project workspace files: %w", err)
	}
	return nil
}

// Init initialises a go workspace in targetDir.
func Init(targetDir string, executor execute.Executor) error {
	if err := Remove(targetDir); err != nil {
		return err
	}
	init := exec.Command("go", "work", "init")
	return executor.Errors(init, targetDir, "initialising workspace")
}

// Use adds all modules in path relative to projectDir to the go.work file.
// Creates a workspace in projectDir if none exist.
func Use(projectDir string, path string, executor execute.Executor) error {
	if err := initIfNoWorkspace(projectDir, executor); err != nil {
		return fmt.Errorf("preparing workspace: %w", err)
	}
	use := exec.Command("go", "work", "use", "-r", path)
	return executor.Errors(use, projectDir, "updating workspace modules")
}

func initIfNoWorkspace(projectDir string, executor execute.Executor) error {
	if !file.Exists(filepath.Join(projectDir, "go.work")) {
		if err := Init(projectDir, executor); err != nil {
			return err
		}
	}
	return nil
}
