# name of output binary
BINARY_NAME=rsync-daemon

# latest git commit hash
LATEST_COMMIT_HASH=$(shell git rev-parse HEAD)

# go commands and variables
GO=go
GOB=$(GO) build
GOM=$(GO) mod
GO_TEST=$(GO) test

# git commands
GIT=git

# environment variables related to
# cross-compilation.
GOOS_MACOS=darwin
GOOS_LINUX=linux
GOARCH=amd64

# currently installed/running Go version (full and minor)
GOVERSION=$(shell go version | grep -Eo '[1-2]\.[[:digit:]]{1,3}\.[[:digit:]]{0,3}')
MINORVER=$(shell echo $(GOVERSION) | awk '{ split($$0, array, ".") } {print array[2]}')

GOPATH=$(shell go env GOPATH)
GOPATH_BIN="$(GOPATH)/bin"

# Color code definitions
# Note: everything is bold.
COLOR_FMT   = \033[1;38;5;$(1)m
GREEN       = $(call COLOR_FMT,70)
BLUE        = $(call COLOR_FMT,27)
RED         = $(call COLOR_FMT,124)
LIGHT_BLUE  = $(call COLOR_FMT,32)
MAGENTA     = $(call COLOR_FMT,128)
RESET_COLOR =\033[0m

COLOR_ECHO = @echo "$(1)$(2)$(RESET_COLOR)"

default: help

# build a vendored binary for macos
macos: BUILD_OS = $(GOOS_MACOS)
macos: build

# build a vendored binary for linx
linux: BUILD_OS = $(GOOS_LINUX)
linux: build

godocs:
	$(call COLOR_ECHO,$(LIGHT_BLUE),"godoc server running on http://localhost:9000")
	@godoc -http=":9000"

# test tecipes

test: check-rsync
	@env TEST_CONFIG_PATH="$(shell echo "$$PWD")/sample_config.yaml" \
	$(GO_TEST) ./...

check-rsync:
ifeq "$(shell which rsync)" ''
	$(call COLOR_ECHO,$(RED),"rsync needs to be installed for $(BINARY_NAME) to be able to do anything")
	@exit 5
else
	$(call COLOR_ECHO,$(GREEN),"rsync is installed!")
	@exit 0
endif

# go mod recipes

mod-clean: ## Remove all the Go mod cache
	@$(GO) clean -modcache

mod-tidy:
	@$(GOM) tidy

mod-vendor: mod-tidy ## run tidy & vendor
	@$(GOM) vendor

# build recipes

build:
	$(call COLOR_ECHO,$(GREEN),"building $(BINARY_NAME) using GOOS=$(BUILD_OS)...")
	@env GOOS=$(BUILD_OS) GOARCH=$(GOARCH) $(GOB) ./cmd/$(BINARY_NAME)

clean:
	$(call COLOR_ECHO,$(GREEN),"cleaning up $(BINARY_NAME)...")
	@rm -rf $(BINARY_NAME)

help: ## This help dialog.
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
	for help_line in $${help_lines[@]}; do \
		IFS=$$'#' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf "%-30s %s\n" $$help_command $$help_info ; \
	done