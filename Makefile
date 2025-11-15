.DEFAULT_GOAL:=help
SHELL:=/bin/sh
GOPATH := $(shell go env GOPATH)

help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

########################################################################################################################
##@ Common
########################################################################################################################

.PHONY: build
build: ## Builds the application for production
	go build -ldflags="-w -s" -v -o ./bin/echo-swagger .

.PHONY: clean
clean: ## Runs mod tidy
	go mod tidy

.PHONY: update
update: ## Update go modules
	go get -t -u ./...

.PHONY: swag
swag: # Install swaggo module and updates all swagger files
	swag init -g ./examples/simple/main.go --parseInternal --outputTypes "go,yaml"
	swag fmt

########################################################################################################################
##@ Run
########################################################################################################################

.PHONY: run
run: ## Execute the application locally
	go run -race .


.PHONY: test
test: # Runs all the tests in the application and returns if they passed or failed, along with a coverage percentage
	go install github.com/mfridman/tparse@main | go mod tidy
	PROFILE=local go test -parallel 10 -json -cover ./... | tparse -all -pass -trimpath=github.com/BrunoKrugel/echo-swagger/

########################################################################################################################
##@ Code Style
########################################################################################################################

.PHONY: format
format: ## Format code and organize imports
	goimports -w .
	go fmt ./...
	fieldalignment -fix ./...

.PHONY: lint
lint: ## Runs golangci-lint
	golangci-lint run

