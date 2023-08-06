# Simple REST API Server

[![Test and build the app](https://github.com/m-kuzmin/simple-rest-api/actions/workflows/golang-ci.yml/badge.svg?branch=main)](https://github.com/m-kuzmin/simple-rest-api/actions/workflows/golang-ci.yml)

# Features

- An endpoint to create user(s) by uploading a CSV file
- [ ] An endpoint to search the users database

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
- TDD: CI (GitHub Actions) for automated builds and checks
- Docker
- REST docs using Swagger UI (`/swaggerui/index.html`)

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

# Commands for development

## Testing the code locally

The Makefile contains a help command that prints out a list of commands. You can add a comment above any command to
include it in the help message. Only the first comment line directly above the command will be printed. Here is an
example:

```makefile
# Help is awesome, but awk is complicated. Also this line isnt part of `make help`.
# Print the help message
help:
```
```shell
$ make help
Malefile help:

help   Print the help message
```

Most of the time you would want to run `make ci` which runs the tests and the linter. However the tests also need a
postgres database so you should start it first:

```shell
docker compose up -d postgres
make ci
```

The `make ci` command performs checks that the remote CI on GitHub would do, just does them locally. When adding a new
check, add it to GitHub Actions workflow first and then to the makefile to make sure GitHub actions contain all
neccessary tests.

# Explanation of the build process

The application has only 1 build script - Dockerfile. There are no other ways to build or run this app. This is because
the build process includes generating swagger docs and the app also needs to have a PostgreSQL database. It is easier to
maintain one build script (Dockerfile + docker-compose.yml) than doing it in a Makefile *and* in docker. The makefile is
just for emulating the CI locally.
