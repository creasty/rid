.DEFAULT_GOAL := all

SHELL := /bin/bash -eu -o pipefail

NAME     := rid
VERSION  := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)

GO_BUILD_FLAGS := -v -ldflags="-s -w -X \"github.com/creasty/rid/cli.Version=$(VERSION)\" -X \"github.com/creasty/rid/cli.Revision=$(REVISION)\" -extldflags \"-static\""
GO_TEST_FLAGS  := -v

PACKAGE_DIRS := $(shell go list ./... 2> /dev/null | grep -v /vendor/)
SRC_FILES    := $(shell find . -name '*.go' -not -path './vendor/*')


#  bin
#-----------------------------------------------
bin/$(NAME): $(SRC_FILES)
	@GOOS=darwin GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o bin/$(NAME)


#  Tasks
#-----------------------------------------------
all: bin/$(NAME)

.PHONY: ci-build
ci-build:
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

.PHONY: ci-test
ci-test: lint
	@go test  $(PACKAGE_DIRS)
	@echo > coverage.txt
	@for d in $(PACKAGE_DIRS); do \
		go test -coverprofile=profile.out -covermode=atomic -race -v $$d; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi; \
	done

.PHONY: release
release:
	git tag v$(VERSION)
	git push origin v$(VERSION)

.PHONY: dist
dist:
	@cd dist \
		&& find * -type d -exec cp ../LICENSE {} \; \
		&& find * -type d -exec cp ../README.md {} \; \
		&& find * -type d -exec tar -zcf $(NAME)-{}.tar.gz {} \; \
		&& find * -type d -exec zip -r $(NAME)-{}.zip {} \; \
		&& cd ..
