package module

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vision-cli/vision/common/execute"
	"github.com/vision-cli/vision/common/file"
)

const (
	modFile   = "go.mod"
	modPrefix = "module"
)

// Remove removes any go.mod or go.sum files in moduleDir.
func Remove(moduleDir string) error {
	if err := file.RemoveNamed(moduleDir, "go.mod", "go.sum"); err != nil {
		return fmt.Errorf("removing existing go module files: %w", err)
	}
	return nil
}

// Init initialises a go module in targetDir with moduleName.
// Removes any existing go module.
func Init(targetDir string, moduleName string, executor execute.Executor) error {
	if err := Remove(targetDir); err != nil {
		return err
	}
	init := exec.Command("go", "mod", "init", moduleName)
	return executor.Errors(init, targetDir, "initialising module")
}

// Tidy tidies module dependencies.
func Tidy(moduleDir string, executor execute.Executor) error {
	tidy := exec.Command("go", "mod", "tidy")
	return executor.Errors(tidy, moduleDir, "finding required module dependencies")
}

// Name returns the module name found in moduleDir/go.mod.
func Name(moduleDir string) (string, error) {
	modPath := filepath.Join(moduleDir, modFile)

	lines, err := file.ToLines(modPath)
	if err != nil {
		return "", fmt.Errorf("reading mod file in %s: %w", moduleDir, err)
	}

	for _, line := range lines {
		if strings.HasPrefix(line, modPrefix) {
			return strings.TrimSpace(strings.Replace(line, modPrefix, "", 1)), nil
		}
	}

	return "", fmt.Errorf("module string not found in %s", modPath)
}

// Rename renames the module in moduleDir/go.mod to newModuleName.
// TODO: replace all references to it in other modules
func Rename(moduleDir string, newModuleName string) error {
	modPath := filepath.Join(moduleDir, modFile)

	lines, err := file.ToLines(modPath)
	if err != nil {
		return fmt.Errorf("reading mod file in %s: %w", moduleDir, err)
	}

	lines[0] = fmt.Sprintf("%s %s", modPrefix, newModuleName)
	if err = file.FromLines(modPath, lines); err != nil {
		return fmt.Errorf("writing new lines to mod file in %s: %w", moduleDir, err)
	}

	return nil
}

// Replace adds a replace directive for serviceMod using replacement
func Replace(moduleDir string, serviceMod string, replacement string, executor execute.Executor) error {
	edit := exec.Command("go", "mod", "edit", "-replace", //nolint:gosec //serviceMod path is cleaned
		fmt.Sprintf("%s=%s", serviceMod, replacement))
	return executor.Errors(edit, moduleDir, fmt.Sprintf("replace for %s", serviceMod))
}
