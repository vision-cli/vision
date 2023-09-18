
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

## Design

The vision cli is really a wrapper around vision plugins. The cli's purpose is to provide helper functions around managing plugins. The plugins will do the code generation.

Ideally projects are created using vision and services are added using vision. This allows vision to manage a project lifecycle and upgrades via its plugins.

Plugins can be made with vision, and should adhere to the vision plugin API.

## Best practises