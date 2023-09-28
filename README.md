
<h1 align="center"><a href="https://en.wikipedia.org/wiki/Vision_(Marvel_Cinematic_Universe">Vision</a></h1>

<p align="center">
  <img width="150" src="./docs/images/vision-3d.jpg" />
</p>

## Roadmap

TODO:

- [ ] Brew install
- [ ] Update all plugins
- [ ] Define vision plugin API

## Install

From source:

Requires [go](https://go.dev/dl/) to be installed

```bash
go install github.com/vision-cli/vision@latest
```

Build from source:

```bash
git clone github.com/vision-cli/vision
go build
# add binary to path
```

## Purpose

Vision is a developer productivity tool which uses plugins to scaffold code templates.
It is intended to make the developer experience better by automating a lot of repetitive
boilerplate.

## Usage

### Initialising

To initialise vision and create a vision.json config file where the project name is the same as the directory name:

```bash
vision init 
```

However, if you want to specify a project name, you can use the -p flag.

```bash
vision init -p [projectname]
```

If you want to create a vision config file inside of a new directory, you can run vision init like so:

```bash
vision init [dirName] -p [projectName]
```

### Creating a plugin


```bash
vision create
```

### Plugin information

To list all of the available plugins:

```bash
vision plugin list
```

If you want to omit all faulty plugins and only show working plugins, you can use the -w flag.

```bash
vision plugin list -w
```

To check plugin health and list reasons why plugins are faulty:

```bash
vision doctor
```

<!-- Add more commands when commands have been created -->

## Design

The vision cli is really a wrapper around vision plugins. The cli's purpose is to provide helper functions around managing plugins. The plugins will do the code generation.

Ideally projects are created using vision and services are added using vision. This allows vision to manage a project lifecycle and upgrades via its plugins.

Plugins can be made with vision, and should adhere to the vision plugin API.

## Best practises

When creating plugins to use with vision, these commands must always be included:

```bash

info

version

init

generate

```

See vision-plugin-sample-v1 for a complete template of what a plugin should look like.
