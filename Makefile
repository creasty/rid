.DEFAULT_GOAL := build

SHELL := /bin/bash -eu -o pipefail

NAME     := dor
VERSION  := 0.0.1
REVISION := $(shell git rev-parse --short HEAD)

GO_BUILD_FLAGS := -v -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""
GO_TEST_FLAGS  := -v -coverprofile=coverage.txt -covermode=atomic

PACKAGE_DIRS := $(shell go list ./... 2> /dev/null | grep -v /vendor/)
SRC_FILES    := $(shell find . -name '*.go' -not -path './vendor/*')


#  Tasks
#-----------------------------------------------
.PHONY: build
build:
	@for os in darwin linux; do \
		for arch in amd64 386; do \
			echo "==> Build $$os $$arch"; \
			GOOS=$$os GOARCH=$$arch go build $(GO_BUILD_FLAGS) \
				-o dist/$$os-$$arch/$(NAME); \
		done; \
	done

.PHONY: clean
clean:
	@rm -rf dist/*

.PHONY: lint
lint:
	@gofmt -e -d -s $(SRC_FILES) | awk '{ e = 1; print $0 } END { if (e) exit(1) }'
	@echo $(SRC_FILES) | xargs -n1 golint -set_exit_status
	@go vet $(PACKAGE_DIRS)

.PHONY: test
test: lint
	@go test $(GO_TEST_FLAGS) $(PACKAGE_DIRS)

.PHONY: release
release:
	git tag $(VERSION)
	git push origin $(VERSION)
