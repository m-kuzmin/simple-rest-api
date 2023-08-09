# Simple REST API Server

[![CI](https://github.com/m-kuzmin/simple-rest-api/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/m-kuzmin/simple-rest-api/actions/workflows/ci.yml)

# Features

**Endpoints:**
- Adding users to the database
- Searching the database by column(s) (partial substring match)
- Swagger UI (default /swaggerui/index.html, can be changed in the config)

A User has the following fields:

- ID
- Name
- Phone number
- Country
- City

# Tech stack

- [Gin](https://github.com/gin-gonic/gin) router
- PostgreSQL Database
- Written in Golang
- TDD: CI (GitHub Actions) for automated testing
- Docker (docker compose)
- Swagger UI for REST API documentation (`/swaggerui/index.html`)

# Running the server

You will need docker installed.

```shell
docker compose up --build -d
```

- `--build` Makes sure the images for all services are (re-)built
- `-d` Makes the containers run in detached mode.

To stop the server run

```shell
docker compose down
```

## Configuration

The `docker-compose.yml` file mounts a volume to the REST API server container. This volume is mapped to `./config/`,
which contains a `config.yml`. Edit this file. Note that 90% of the config must match the docker-compose file because of
things like postgres address, user, password, etc.

# Commands for development

## Testing the code locally

The makefile serves as a replica of the GitHub Actions process. The default `ci` rule runs tests and linters. You can
run just one step of the CI process with e.g. `make lint`. All documented rules can be viewed with `make help`. 
Documented simply means the line above the rule starts with `# Some doc string`. All rules *should* have a doc comment,
but it might be missing for some of them.

If it is neccessary to test a new feature of the app the easiest way to do so is to write some tests in Go. Otherwise
you can use Swagger UI. But the issue with doing all testing in swagger is that you will likely be the last person to
ever test this functionality, compared to automated tests which will stay there until the requirements change.

# Features

**Endpoints:**
- Adding users to the database
- Searching the database by column(s) (partial substring match)
- Swagger UI (default /swaggerui/index.html, can be changed in the config)

A User has the following fields:

- ID
- Name
- Phone number
- Country
- City

# Tech stack

- [Gin](https://github.com/gin-gonic/gin) router
- PostgreSQL Database
- Written in Golang
- TDD: CI (GitHub Actions) for automated testing
- Docker (docker compose)
- Swagger UI for REST API documentation (`/swaggerui/index.html`)

# Running the server

You will need docker installed.

```shell
docker compose up --build -d
```

- `--build` Makes sure the images for all services are (re-)built
- `-d` Makes the containers run in detached mode.

To stop the server run

```shell
docker compose down
```

## Configuration

The `docker-compose.yml` file mounts a volume to the REST API server container. This volume is mapped to `./config/`,
which contains a `config.yml`. Edit this file. Note that 90% of the config must match the docker-compose file because of
things like postgres address, user, password, etc.

# Commands for development

## Testing the code locally

The makefile serves as a replica of the GitHub Actions process. The default `ci` rule runs tests and linters. You can
run just one step of the CI process with e.g. `make lint`. All documented rules can be viewed with `make help`. 
Documented simply means the line above the rule starts with `# Some doc string`. All rules *should* have a doc comment,
but it might be missing for some of them.

If it is neccessary to test a new feature of the app the easiest way to do so is to write some tests in Go. Otherwise
you can use Swagger UI. But the issue with doing all testing in swagger is that you will likely be the last person to
ever test this functionality, compared to automated tests which will stay there until the requirements change.
