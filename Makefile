# This makefile is for developers. The application can be started with:
# ```shell
# docker compose up --build -d
# ```

.PHONY: default ci test lint

# (default) Run the actions that the CI would run
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
