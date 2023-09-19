package cmd

import (
	"errors"
	"log"

	"github.com/vision-cli/vision/cli"
	"github.com/vision-cli/vision/cmd/config"
	"github.com/vision-cli/vision/common/execute"
	cc "github.com/vision-cli/vision/common/plugins"
	rp "github.com/vision-cli/vision/remote-plugins"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(config.RootCmd)
	osExecutor := execute.NewOsExecutor()
	p, err := cc.GetPlugins(osExecutor)
	if err != nil {
		cli.Warningf("cannot get plugins: %v", err)
	}
	for _, pl := range p {
		cobraCmd, err := rp.GetCobraCommand(pl, osExecutor)
		if err != nil {
			cli.Warningf("cannot get cobra command %s: %v", pl.Name, err)
		}
		rootCmd.AddCommand(cobraCmd)
	}
}

var rootCmd = &cobra.Command{
	Use:   "vision",
	Short: "A developer productivity tool",
	Long:  `Vision is tool to create microservice platforms and microservice scaffolding code`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return checkTools(execute.NewOsExecutor())
	},
	Example: `You need to create a seed project in the cloud you want before using Vision for the first time.
The seed project will be used to store terraform state and hold the container registry for your microservices.
	
Run the following command to create a new project

	vision project create myproject -r github.com/myorg/myproject -g gcr.io/myproject

This will create a folder called myproject with the standard vision folder structure and a default config file.
There is a powerful option to create a project from a template model using

	vision project create -t <template-name>

See examples folder for example template files
Once you have created a project, navigate to the project folder and create a microservice using

	vision service create myservice

This will create a folder called myservice with the standard vision folder structure for a microservice.
Create a microservice platform for a cloud provider, for example creating an Azure platform using

	vision infra create azure -d standalone-graphql
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cli.Fatalf(err.Error())
	}
	log.Println("root command executes")
}

func checkTools(e execute.Executor) error {
	if !e.CommandExists("go") {
		return errors.New("go is not installed")
	}
	if !e.CommandExists("protoc") {
		cli.Warningf("The protoc cli is not installed. You will need this to build the service. See https://grpc.io/docs/protoc-installation/ for installation instructions.")
	}
	if !e.CommandExists("dapr") {
		cli.Warningf("The dapr cli is not installed. You will need this to run your service locally. See https://docs.dapr.io for installation instructions.")
	}
	if !e.CommandExists("docker") {
		cli.Warningf("The docker cli is not installed. You will need this to build infrastructure. See (https://www.docker.com) for installation instructions.")
	}
	if !e.CommandExists("az") {
		cli.Warningf("The az cli is not installed. You will need this to build Azure infrastructure. See (https://learn.microsoft.com/en-us/cli/azure/install-azure-cli) for installation instructions.")
	}
	if !e.CommandExists("aws") {
		cli.Warningf("The aws cli is not installed. You will need this to build Aws infrastructure. See (https://aws.amazon.com/cli/) for installation instructions.")
	}
	if !e.CommandExists("gcloud") {
		cli.Warningf("The gcloud cli is not installed. You will need this to build Gcp infrastructure. See (https://cloud.google.com/sdk/gcloud) for installation instructions.")
	}
	if !e.CommandExists("terraform") {
		cli.Warningf("The terraform cli is not installed. You will need this to build infrastructure. See (https://www.terraform.io) for installation instructions.")
	}
	return nil
}
