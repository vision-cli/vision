# NHS Scotland SCI Gateway

## Development

### Getting started

#### Running the app

To run the app, run the following command from the root of the project:

```bash
make run
```

**Note:** If you get an error about connecting to a database, see the [Prerequisites](#prerequisites) section below.

### Prerequisites

#### Postgres

The app requires a postgres database to be running. The easiest way to do this is to use docker-compose. Ensure you have docker installed and run the following command from the root of the project:

```bash
cd docker/postgres
docker-compose up -d
```

#### (Optional) Auto reload

Reloading the go build each time a file changes to ease front end development can be a nice quality of life feature. This can be achieved with [air](https://github.com/cosmtrek/air)

The included [air.toml](air.toml) file is configured to watch for changes. To use it, install air with:

```bash
go install github.com/cosmtrek/air@latest
```

Then run the following command from the root of the project:

```bash
air
```

#### (Optional) Prettier formatting

We aim to avoid node as much as possible so we do not have any part of the node ecosystem for code. However, we do use prettier for formatting. To use it, install prettier with:

```bash
npm install
```

The vscode config and prettier config should take care of the rest.
