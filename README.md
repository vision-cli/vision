
<h1 align="center"><a href="https://en.wikipedia.org/wiki/Vision_(Marvel_Cinematic_Universe">Vision</a></h1>

<p align="center">
  <img width="150" src="./images/vision-3d.jpg" />
</p>

## Purpose

Vision is an evolution of a developer productivity tool called Jarvis that was built and used by a full stack engineering
team over many years. The tools are intended to
write scaffolding code using sane defaults and best practices for a developer to then modify and extend.
Vision is not intended to replace a developer, nor will it make the task of coding less complex (for example
the claim of low-code platforms), it is intended to make the developer experience better by automating a lot of repettive
scaffoding code writing.

The vision cli is really a wrapper around vision plugins. The cli's purpose is to provide help to the user, ensure the user is
in the correct folder, create and manage config and use command line flags to override that config before calling plugins and
passing that config to them. The plugins will do the code generation.

Vision plugins use templates and Asbtract Syntax Trees (<https://en.wikipedia.org/wiki/Abstract_syntax_tree>) to write and
manipulate code. Vision expects a project's folder and file structure to be in accordance with its standard.
Ideally projects are created using vision and services are added using vision.

## Dependencies

Vision requires golang (<https://go.dev>) to be installed.
Vision plugins have other dependencies including

- docker (<https://www.docker.com>)
- dapr (<https://dapr.io>)
- terraform (<https://www.terraform.io>)
- azure cli (<https://learn.microsoft.com/en-us/cli/azure/install-azure-cli>)
- aws cli (<https://aws.amazon.com/cli/>)
- gcloud cli (<https://cloud.google.com/sdk/gcloud>)

Install the grpc tools by running:

```
go install \
  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
  google.golang.org/protobuf/cmd/protoc-gen-go \
  google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

This will place four binaries in your $GOBIN;

    protoc-gen-grpc-gateway
    protoc-gen-openapiv2
    protoc-gen-go
    protoc-gen-go-grpc

Make sure that your $GOBIN is in your $PATH

## Usage

Install vision with

```
go install github.com/vision-cli/vision
```

(Ensure $GOBIN is on the system $PATH)

When you run vision for the first time it will download a standard set of plugins for creating
project, services, graphql service, gateway service and infra.

### Basic Project Template

Get started with

```
vision project create <project name>
```

to generate a project directory with configuration.
(If no projectName is specified, configuration will be created in the current working directory)

The project will contain a basic project setup. Change to the directory vision has just created (your project name) and run the
following command to create a basic service template

```
vision service create <service name>
```

### Advanced Projects Templates

Create a yaml file with a description of your project's model. The schema and an example can be found in the example folder
(example_model.yml). Run the command below to create the project including

- golang models, updated golang command, server and resolvers to handle messages and persistence
- protobuf services and messages
- graphql model, queries and mutations

```
vision project create -t <path to your model.yml>
```

### Other Useful Commands

| Command                                                    |                          Description                           |
| :--------------------------------------------------------- | :------------------------------------------------------------: |
| `vision service endpoints <serviceName> -n <namespace>`    |    generate rest endpoints in `grpc-gateway-api-config.yml`    |
| `vision service update <serviceName> --rn <newName>`       | rename a service (default namespace, -n for other, -A for all) |
| `vision service update -n <namespace> --mv <newNamespace>` |  move services to another namespace (can specify serviceName)  |
| `vision project docs`                                      |         generate a table of contents and services list         |

### Gateway Generation

Configure rest endpoints for rpc methods in `grpc-gateway-api-config.yml`.
Within your grpc service, run `make proxy` to create gateway stubs.
This will generate the file {namespace}\_{serviceName}.pb.gw.go in the proto/ directory.
Then generate your grpc gateway service to expose grpc over rest.

```
vision gateway create <grpcGatewayServiceName> -n <namespace>
```

### Principles of Commands

- All commands should be callable with mininal parameters, and the command will behave interactuvely
- All requred parameters for a command should also be "providable" as parameters on the command line
- There should always be a silent option to select defaults answers to interactive questions
- All commands should be well documented with examples
- Commands will always update the json config file with the latest settings used
