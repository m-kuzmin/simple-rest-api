# Simple REST API Server

[![Test and build the app](https://github.com/m-kuzmin/simple-rest-api/actions/workflows/golang-ci.yml/badge.svg?branch=main)](https://github.com/m-kuzmin/simple-rest-api/actions/workflows/golang-ci.yml)

# Features

- [ ] An endpoint to create user(s) by uploading a CSV file
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
- [ ] Docker
- [ ] REST docs using Swagger

# Running the server

You will need docker installed.

```shell
docker compose up
```

To stop the server run

```shell
docker compose down
```

# Commands for development

## Checking the code locally

The Makefile contains a help command that prints out a list of commands. You can add a comment above any command to
include it in the help message. Only the first comment line directly above the command will be printed.

Most of the time you would want to run either:

```shell
make ci
# or
make
```

The `ci` command performs checks that the remote CI on GitHub would do, just does them locally. When adding a new check,
add it to actions first and then to the makefile to make sure GitHub actions contain all checks neccessary.
