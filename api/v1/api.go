package api

// A plugin must accept a json serialized PluginRequest from stdin and return a json serialized PluginResponse.
// The commands present in the v1 api are listed in the Commands constants below.
// A plugin must return a json serialized response based on the command sent. The response types are indicated below.
// If an error occurs during the execution of the plugin, the plugin must return a PluginResponse with the
// 'Result' set to an empty string and the error message in 'Error'.

// Commands that the plugin must support
const (
	CommandUsage  = "usage"
	CommandConfig = "config"
	CommandRun    = "run"
)

// PluginRequest is the structure that is passed to the plugin from the CLI
type PluginRequest struct {
	Command      string             // v1 plugins must support commands listed in the Commands constants above
	Args         []string           // blank if the command is 'usage' or 'config', otherwise the sub commands e.g. create, delete, etc if the command is 'run'
	Flags        []PluginFlag       // blank if the command is 'usage' or 'config', otherwise the flags structure provided at the command line
	Placeholders PluginPlaceholders // blank if the command is 'usage' or 'config', otherwise the placeholders structure built from the config file and flags
}

// PluginResponse is the default structure that is returned from the plugin to the CLI
type PluginResponse struct {
	Result string // Result of the command displayed to the user
	Error  string // Fatal error message displayed to the user
}

// PluginUsageResponse is the structure that is returned from the plugin to the CLI when the command is 'usage'
type PluginUsageResponse struct {
	Version        string       // API version supported (the package name above)
	Use            string       // The command name
	Short          string       // Short description shown in the 'usage' output
	Long           string       // Long description shown in the 'usage' output
	Example        string       // Example shown in the 'usage' output
	Subcommands    []string     // Subcommands shown in the 'usage' output
	Flags          []PluginFlag // Flags required in addition to the standard ConfigFlags
	RequiresConfig bool         // Whether the plugin requires a config file
}

type PluginFlag struct {
	Name      string // name of the flag
	Shorthand string // shorthand name of the flag. Must be unique
	Value     string // default value of the flag
	Usage     string // usage description of the flag
}

// PluginConfigResponse is the structure that is returned from the plugin to the CLI when the command is 'config'
type PluginConfigResponse struct {
	Defaults []PluginConfigItem // the default config items that the CLI will write to the config file
}

type PluginConfigItem struct {
	Key     string // key of the config item
	Default string // default value of the config item
}

// The placeholders structure is built from the config file, flags, and best practice conventions that
// is passed to the plugin from the CLI. The placeholders structure is used to build the templates
type PluginPlaceholders struct {
	// ----------------------- project --------------------------

	// Project Root - this is the root project folder where services, packages, etc can be found
	// relative to where the command has been run, for example if the command is run from the
	// the project folder itself e.g. user/demo, the Project Root will be "." whereas if
	// the command is run from the "user" folder then the Project Root will be "demo"
	// The go.work file should be placed here
	ProjectRoot string

	// Name of the the project regardless of the project root e.g. myproject
	ProjectName string

	// Name of the the project directory e.g. my-project
	ProjectDirectory string

	// The fully qualified namespace of the project e.g. github.com/mycompany/myproject
	// This is constructed from <Remote>/<ProjectName>
	ProjectFqn string

	// the container registry to use for the project. This can be overridden with the -g --registry flag
	Registry string

	// Remote directory used for eg github.com/mycompany. This can be overridden with the -r --remote flag
	Remote string

	// default branch to use for the module, e.g. master. This can be overridden with the -b --branch flag
	Branch string

	// major version of the service api, e.g. v1. This can be overridden with the -v --version flag
	Version string

	// A unique string used to generate unique names for resources e.g. myproject-agshrt
	UniqueStr string

	// ----------------------- services --------------------------

	// The fully qualified namespace of default project services e.g. github.com/mycompany/myproject/services/default
	// This is constructed from <ProjectFqn>/<ServicesDirectory>/<DefautlServiceNamespace>
	ServicesFqn string

	// the services directory e.g. services
	ServicesDirectory string

	// gateway service name e.g. gateway
	GatewayServiceName string

	// fully qualified name of the gateway service e.g. github.com/mycompany/myproject/services/default/gateway
	GatewayFqn string

	// graphql service name e.g. graphql
	GraphqlServiceName string

	// fully qualified name of the graphql service e.g. github.com/mycompany/myproject/services/default/graphql
	GraphqlFqn string

	// the fullt qualified project libs namespace name e.g. github.com/mycompany/myproject/libs
	LibsFqn string

	// the full libs directory e.g. ./libs
	// this is formed from the <ProjectRoot>/<LibsDirectory>
	LibsDirectory string

	// ----------------------- service --------------------------

	// The namespace of the service e.g. default
	ServiceNamespace string

	// The service name and version e.g. default.v1
	ServiceVersionedNamespace string

	// name of the service to create or operate on e.g. testservice
	ServiceName string

	// The fully qualified namespace of the service e.g. github.com/mycompany/myproject/services/default/testservice
	// This is constructed from <ProjectFqn>/<ServiceNamespace>/<ServiceName>
	ServiceFqn string

	// the directory name of the service to create or operate on e.g. ./services/default/testservice
	ServiceDirectory string

	// ----------------------- infra --------------------------

	// the infra directory e.g. infra
	InfraDirectory string

	// ----------------------- messaging --------------------------

	// protocol buffer package name, e.g. namespace.v1_testservice.v1
	ProtoPackage string

	// ----------------------- graphql --------------------------

	// deployment type, e.g. standalone-graphql, standalone-gateway
	Deployment string
}
