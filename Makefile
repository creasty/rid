.DEFAULT_GOAL := all

SHELL := /bin/bash -eu -o pipefail
ROOT_DIR := $(shell pwd)

NAME     := rid
VERSION  := 0.1.0
REVISION := $(shell git rev-parse --short HEAD)

REPO := github.com/creasty/rid

PACKAGE_DIRS := $(shell go list ./... 2> /dev/null | grep -v /vendor/)
SRC_FILES := $(shell git ls-files --cached --others --exclude-standard | grep -E "\.go$$")

CMD_DIR := ./cmd
BIN_DIR := ./bin
DST_DIR := ./dist

XC_ARCH := 386 amd64
XC_OS := darwin linux

#  Flags
#-----------------------------------------------
GO_BUILD_FLAGS := -v
GO_TEST_FLAGS := -v
GO_LDFLAGS := \
	-s -w \
	-X '$(REPO)/cmd.Version=$(VERSION)' \
	-X '$(REPO)/cmd.Revision=$(REVISION)' \
	-extldflags '-static'
GO_COVER_FLAGS := \
	-v \
	-coverpkg $(shell echo $(PACKAGE_DIRS) | tr ' ' ',') \
	-coverprofile coverage.txt \
	-covermode atomic \
	-race

#  Godep
#-----------------------------------------------
DEP_VENDOR_PATH := $(ROOT_DIR)/vendor
DEP_BIN_PATH := $(DEP_VENDOR_PATH)/.bin
DEP_CMDS := \
	github.com/mitchellh/gox \
	github.com/golang/mock/mockgen

DEP_BINS := $(addprefix $(DEP_BIN_PATH)/,$(notdir $(DEP_CMDS)))

define dep-bin-tmpl
$(DEP_BIN_PATH)/$(notdir $(1)): dep
	@echo "Installing $(1)"
	@cd $(DEP_VENDOR_PATH)/$(1) && GOBIN="$(DEP_BIN_PATH)" go install .
endef

$(foreach src,$(DEP_CMDS),$(eval $(call dep-bin-tmpl,$(src))))

#  Bin
#-----------------------------------------------
$(BIN_DIR)/$(NAME): $(SRC_FILES) gen
	@go build \
		$(GO_BUILD_FLAGS) \
		-ldflags "$(GO_LDFLAGS)" \
		-o $(BIN_DIR)/$(NAME) \
		$(CMD_DIR)

#  Tasks
#-----------------------------------------------
all: $(BIN_DIR)/$(NAME)

.PHONY: setup
setup: dep $(DEP_BINS)

.PHONY: dep
dep: Gopkg.toml Gopkg.lock
	@dep ensure -v

.PHONY: gen
gen: $(SRC_FILES)
	@PATH=$(DEP_BIN_PATH):$$PATH go generate ./...

.PHONY: lint
lint:
	@gofmt -e -d -s $(SRC_FILES) | awk '{ e = 1; print $0 } END { if (e) exit(1) }'
	@echo $(SRC_FILES) | xargs -n1 golint -set_exit_status
	@go vet $(PACKAGE_DIRS)

.PHONY: test
test: gen lint
	@go test $(GO_TEST_FLAGS) ./...

.PHONY: ci-test
ci-test: gen lint
	@echo > coverage.txt
	@go test $(GO_TEST_FLAGS) $(GO_COVER_FLAGS) ./...

.PHONY: release
release:
	git tag v$(VERSION)
	git push origin v$(VERSION)

.PHONY: dist
dist: gen
	@PATH=$(DEP_BIN_PATH):$$PATH gox \
		-ldflags="$(GO_LDFLAGS)" \
		-os="$(XC_OS)" \
		-arch="$(XC_ARCH)" \
		-output="$(DST_DIR)/$(NAME)_{{.OS}}_{{.Arch}}" \
		$(CMD_DIR)

.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)/*
	@rm -rf $(DST_DIR)/*
