.DEFAULT_GOAL := build

SHELL := /bin/bash -eu -o pipefail

NAME     := dor
VERSION  := 0.0.1
REVISION := $(shell git rev-parse --short HEAD)

GO_BUILD_FLAGS := -v -ldflags="-s -w -X \"github.com/creasty/dor.Version=$(VERSION)\" -X \"github.com/creasty/dor.Revision=$(REVISION)\" -extldflags \"-static\""
GO_TEST_FLAGS  := -v -race

PACKAGE_DIRS := $(shell go list ./... 2> /dev/null | grep -v /vendor/)
SRC_FILES    := $(shell find . -name '*.go' -not -path './vendor/*')


#  Tasks
#-----------------------------------------------
.PHONY: build
build:
	@for os in darwin linux windows; do \
		for arch in amd64 386; do \
			echo "Build $$os $$arch"; \
			GOOS=$$os GOARCH=$$arch go build $(GO_BUILD_FLAGS) \
				-o dist/$$os-$$arch/$(NAME); \
		done; \
	done

.PHONY: clean
clean:
	@rm -rf dist/*

.PHONY: lint
lint:
	@gofmt -e -d -s $(SRC_FILES) | awk '{ E=1; print $0 } END { if (E) exit(1) }'
	@echo $(SRC_FILES) | xargs -n1 golint -set_exit_status
	@go vet $(PACKAGE_DIRS)

.PHONY: test
test: lint
	@go test $(GO_TEST_FLAGS) $(PACKAGE_DIRS)

.PHONY: release
release:
	git tag $(VERSION)
	git push origin $(VERSION)
