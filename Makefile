SHELL:=/bin/bash

PROJECT_PATH := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
TEST_OUTPUT_PATH := $(PROJECT_PATH)/.output


.PHONY: unit
unit:
	go mod tidy && go test ./... -race -coverprofile .testCoverage.txt

.PHONY: unit-coverage
unit-coverage: unit ## Runs unit tests and generates a html coverage report
	go tool cover -html=.testCoverage.txt -o unit.html
