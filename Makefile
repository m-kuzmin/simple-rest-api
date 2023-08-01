# Executing the application and testing
.PHONY: default run ci test lint

# Test and run the app
default: ci run

# Run the app natively (doesnt build)
run:
	@echo -ne "\n~~~ Starting the application: "
	go run main.go

# Run the actions that the CI would run
ci: test lint
	@echo -ne "\n~~~ Checking: "
	sqlc diff
	@
	@echo -ne "\n~~~ Checking: "
	sqlc vet

# Test the code
test:
	@echo -ne "\n~~~ Running tests: "
	go test -coverprofile /dev/null ./...

# Run the linter
lint:
	@echo -ne "\n~~~ Running linter: "
	golangci-lint run

# File generation tools
.PHONY: swagger

# Generate swagger YAML and JSON files
swagger:
	@echo -ne "\n~~~ Generating Swagger documentation: "
	swag init

.PHONY: help

# Help can only be displayed for rules that contain AZ09.-_ and have a comment immediately above the rule. Only the
# comment text on the line before the rule is displayed.
# Show this help.
help:
	@# https://stackoverflow.com/questions/35730218/how-to-automatically-generate-a-makefile-help-command
	@echo "$(MAKEFILE_LIST) help:"
	@echo
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:].-_]+:/{print substr($$1,1,index($$1,":")),c}1{c=0}' \
		$(MAKEFILE_LIST) | column -s: -t
