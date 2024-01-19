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

### 1. Initialise a Vision project.

`init` can accept a project name with the `-p` flag, where it creates a `vision.json` in the current working directory, and/or it can accept an optional argument that creates a new directory. Not including the `-p` flag assigns the name of the directory as the project name in `vision.json`.

```bash
tonystark:~ vision init [DIR] [-p/--project PROJECTNAME]
```

```bash
tonystark:~ vision init avengers

vision init avengers-assemble \
  -p=avengers \
  -m="github.com/stark-industries/avengers-assemble"
```

This creates a directory called `avengers` and creates a `vision.json` file within `avengers`. It also assigns the `project_name` as `avengers`.

```json
{
  "project_name": "avengers"
}
```

### 2. Initialise plugins to use inside project.

We'll initialise a new Go REST server inside the avengers directory.

```bash
tonystark:~/avengers vision gorest init
```

This adds a template configuration to the `vision.json` file which needs to be changed to your specific data.

### 3. Generate the template files from the plugin

To create a

```bash
tonystark:~/avengers vision plugin init
```

```
- myproject
|__ services
  |__ user_service
    |__  cmd
    |__ internal
|__ vision.json
```

### Visual tutorial/demo

To see an example of how to successfully use a vision plugin, you can watch [this video](https://asciinema.org/a/WJD7PJUvkVyDMzl6oleSGhv6i).

## Design

The vision cli is really a wrapper around vision plugins. The cli's purpose is to provide helper functions around managing plugins. The plugins will do the code generation.

Ideally projects are created using vision and services are added using vision. This allows vision to manage a project lifecycle and upgrades via its plugins.

Plugins can be made with vision, and should adhere to the vision plugin API.

## Best practises
