package run

import (
	"embed"
	"fmt"
	"os/exec"
	"path/filepath"

	api_v1 "github.com/vision-cli/api/v1"

	"github.com/vision-cli/vision/common/execute"
	"github.com/vision-cli/vision/common/file"
	"github.com/vision-cli/vision/common/module"
	"github.com/vision-cli/vision/common/tmpl"
	"github.com/vision-cli/vision/common/workspace"
)

const (
	goTemplateDir = "_templates"
)

//go:embed all:_templates
var templateFiles embed.FS

func Create(p *api_v1.PluginPlaceholders, executor execute.Executor, t tmpl.TmplWriter) error {
	var err error
	if err = tmpl.GenerateFS(templateFiles, goTemplateDir, p.ProjectRoot, p, true, t); err != nil {
		return fmt.Errorf("generating the project structure from the template: %w", err)
	}
	if err = GenerateDocs(p, p.ProjectRoot, false, t); err != nil {
		return fmt.Errorf("generating the project documentation: %w", err)
	}
	if !file.Exists(".git") {
		if err = initGit(p.ProjectRoot, p.Branch, executor); err != nil {
			return fmt.Errorf("initialising the git repository for the project: %w", err)
		}
	}
	persistenceLibRoot := filepath.Join(p.ProjectRoot, p.LibsDirectory, "go", "persistence")
	persistenceFqn := filepath.Join(p.LibsFqn, "go", "persistence")
	err = module.Init(persistenceLibRoot, persistenceFqn, executor)
	if err != nil {
		println(err.Error())
		return err
	}
	err = module.Tidy(persistenceLibRoot, executor)
	if err != nil {
		return err
	}
	if err = createWorkspace(p.ProjectRoot, executor); err != nil {
		return fmt.Errorf("creating project workspace: %w", err)
	}
	if err = createConfig(p, executor); err != nil {
		return fmt.Errorf("creating project workspace: %w", err)
	}
	return nil
}

func createWorkspace(targetDir string, executor execute.Executor) error {
	var err error
	if err = workspace.Init(targetDir, executor); err != nil {
		return fmt.Errorf("creating go workspace file: %w", err)
	}
	println("adding project modules to workspace")
	if err = workspace.Use(targetDir, ".", executor); err != nil {
		return fmt.Errorf("adding all project modules to workspace: %w", err)
	}
	return nil
}

func initGit(targetDir, branch string, executor execute.Executor) error {
	action := fmt.Sprintf("initialising project git repository - default branch: [%s]", branch)

	gitInitCmd := exec.Command("git", "init", "-b", branch)
	err := executor.Errors(gitInitCmd, targetDir, action)
	if err != nil {
		return err // error wrapped with action string
	}

	return nil
}

func createConfig(p *api_v1.PluginPlaceholders, executor execute.Executor) error {
	action := "creating default config"

	gitInitCmd := exec.Command("vision", "config", "create", "--silent", "-r", p.Remote, "-g", p.Registry)
	err := executor.Errors(gitInitCmd, p.ProjectDirectory, action)
	if err != nil {
		return err // error wrapped with action string
	}

	return nil
}
