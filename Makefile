-include .env

# Goodies
V = 0
Q = $(if $(filter 1,$V),,@)
E := 
S := $E $E
M = $(shell printf "\033[34;1m▶\033[0m")
rwildcard = $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))

# Folders
BIN_DIR  ?= $(CURDIR)/bin
DEST_DIR ?= /usr/bin
LOG_DIR  ?= log
TMP_DIR  ?= tmp
COV_DIR  ?= tmp/coverage

# Version, branch, and project
BRANCH    != git symbolic-ref --short HEAD
COMMIT    != git rev-parse --short HEAD
BUILD     := "$(STAMP).$(COMMIT)"
VERSION   != awk '/^var +VERSION +=/{gsub("\"", "", $$4) ; print $$4}' version.go
ifeq ($(VERSION),)
VERSION   != git describe --tags --always --dirty="-dev"
endif
REVISION  ?= 1
PROJECT   != awk '/^const +APP += +/{gsub("\"", "", $$4); print $$4}' version.go
ifeq (${PROJECT},)
PROJECT   != basename "$(PWD)"
endif
PACKAGE   = bitbucket-cli
PACKAGE   ?= $(PROJECT)
PLATFORMS ?= darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 windows-amd64 windows-arm64
export PACKAGE PROJECT VERSION BRANCH COMMIT BUILD REVISION

# Files
GOTESTS   := $(call rwildcard,,*_test.go)
GOFILES   := $(filter-out $(GOTESTS), $(call rwildcard,,*.go))
ASSETS    :=

# Testing
TEST_TIMEOUT  ?= 30
COVERAGE_MODE ?= count
COVERAGE_OUT  := $(COV_DIR)/coverage.out
COVERAGE_XML  := $(COV_DIR)/coverage.xml
COVERAGE_HTML := $(COV_DIR)/index.html

# Tools
GO      ?= go
GOOS    != $(GO) env GOOS
LOGGER   =  bunyan -L -o short
GOBIN    = $(BIN_DIR)
GOLINT  ?= golangci-lint
YOLO     = $(BIN_DIR)/yolo
GOCOV    = $(BIN_DIR)/gocov
GOCOVXML = $(BIN_DIR)/gocov-xml
NFPM     = nfpm
GOMPLATE = gomplate
PANDOC  ?= pandoc
TAR     ?= tar
7ZIP    ?= 7z
ZIP     ?= zip
MOVE    ?= mv
COPY    ?= cp -f

# Flags
#MAKEFLAGS += --silent
# GO
export GOPRIVATE   ?= bitbucket.org/gildas_cherruel/*
export CGO_ENABLED  = 0
ifneq ($(what),)
TEST_ARG := -run '$(what)'
else
TEST_ARG :=
endif

ifeq ($(OS), Windows_NT)
  OSTYPE = windows
  OSARCH = amd64
  include Makefile.windows
else
  OSTYPE != uname -s
  OSARCH != uname -p
  ifeq ($(OSTYPE), Linux)
    OSTYPE = linux
    ifeq ($(OSARCH), x86_64)
      OSARCH = amd64
    else ifeq ($(OSARCH), aarch64)
      OSARCH = arm64
    endif
    include Makefile.linux
  else ifeq ($(OSTYPE), Darwin)
    OSTYPE = darwin
    ifeq ($(OSARCH), x86_64)
      OSARCH = amd64
    else ifeq ($(OSARCH), aarch64)
      OSARCH = arm64
    endif
    include Makefile.linux
  else ifeq ($(OSTYPE),)
    $(error Please use GNU Make 4 at least)
  else
    $(error Unsupported Operating System)
  endif
endif

# Main Recipes
.PHONY: all archive build dep fmt gendoc help install lint logview publish run start stop test version vet watch

help: Makefile; ## Display this help
	@echo
	@echo "$(PROJECT) version $(VERSION) build $(BUILD) in $(BRANCH) branch"
	@echo "Make recipes you can run: "
	@echo
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) |\
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo

all: test build; ## Test and Build the application

gendoc: __gendoc_init__ $(BIN_DIR)/$(PROJECT).pdf; @ ## Generate the PDF documentation

publish: __publish_init__ __publish_binaries__; @ ## Publish the binaries to the Repository

archive: __archive_init__ __archive_all__ __archive_chocolatey__ __archive_debian__ __archive_rpm__ ; @ ## Archive the binaries

build: __build_init__ __build_all__; @ ## Build the application for all platforms

install: $(BIN_DIR)/$(OSTYPE)-$(OSARCH)/$(PROJECT); @ ## Install the application
	$(info $(M) Installing application for $(OSTYPE) on $(OSARCH) in $(DEST_DIR)...)
	$Q install $(BIN_DIR)/$(OSTYPE)-$(OSARCH)/$(PROJECT) $(DEST_DIR)/$(PROJECT)

dep:; $(info $(M) Updating Modules...) @ ## Updates the GO Modules
	$Q $(GO) get -u ./...
	$Q $(GO) mod tidy

lint:;  $(info $(M) Linting application...) @ ## Lint Golang files
	$Q $(GOLINT) run *.go

fmt:; $(info $(M) Formatting the code...) @ ## Format the code following the go-fmt rules
	$Q $(GO) fmt *.go

vet:; $(info $(M) Vetting application...) @ ## Run go vet
	$Q $(GO) vet *.go

run:; $(info $(M) Running application...) @  ## Execute the application
	$Q $(GO) run . | $(LOGGER)

logview:; @ ## Open the project log and follows it
	$Q tail -f $(LOG_DIR)/$(PROJECT).log | $(LOGGER)

clean:; $(info $(M) Cleaning up folders and files...) @ ## Clean up
	$Q rm -rf $(BIN_DIR)  2> /dev/null
	$Q rm -rf $(LOG_DIR)  2> /dev/null
	$Q rm -rf $(TMP_DIR)  2> /dev/null

version:; @ ## Get the version of this project
	@echo $(VERSION)

# Development server (Hot Restart on code changes)
start:; @ ## Run the server and restart it as soon as the code changes
	$Q bash -c "trap '$(MAKE) stop' EXIT; $(MAKE) --no-print-directory watch run='$(MAKE) --no-print-directory __start__'"

restart: stop __start__ ; @ ## Restart the server manually

stop: | $(TMP_DIR); $(info $(M) Stopping $(PROJECT) on $(GOOS)) @ ## Stop the server
	$Q-touch $(TMP_DIR)/$(PROJECT).pid
	$Q-kill `cat $(TMP_DIR)/$(PROJECT).pid` 2> /dev/null || true
	$Q-rm -f $(TMP_DIR)/$(PROJECT).pid

# Tests
TEST_TARGETS := test-default test-bench test-short test-failfast test-race
.PHONY: $(TEST_TARGETS) test tests test-ci
test-bench:    ARGS=-run=__nothing__ -bench=. ## Run the Benchmarks
test-short:    ARGS=-short                    ## Run only the short Unit Tests
test-failfast: ARGS=-failfast                 ## Run the Unit Tests and stop after the first failure
test-race:     ARGS=-race                     ## Run the Unit Tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
test tests: | coverage-tools; $(info $(M) Running $(NAME:%=% )tests...) @ ## Run the Unit Tests (make test what='TestSuite/TestMe')
	$Q mkdir -p $(COV_DIR)
	$Q $(GO) test \
			-timeout $(TEST_TIMEOUT)s \
			-covermode=$(COVERAGE_MODE) \
			-coverprofile=$(COVERAGE_OUT) \
			-v $(ARGS) $(TEST_ARG) .
	$Q $(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_OUT) | $(GOCOVXML) > $(COVERAGE_XML)

test-ci:; @ ## Run the unit tests continuously
	$Q $(MAKE) --no-print-directory watch run="make test"
test-view:; @ ## Open the Coverage results in a web browser
	$Q xdg-open $(COV_DIR)/index.html

# Folder recipes
$(BIN_DIR): ; $(MKDIR)
$(TMP_DIR): ; $(MKDIR)
$(LOG_DIR): ; $(MKDIR)
$(COV_DIR): ; $(MKDIR)

# Documentation recipes
__gendoc_init__:; $(info $(M) Building the documentation...)

$(BIN_DIR)/$(PROJECT).pdf: README.md ; $(info $(M) Generating PDF documentation in $(BIN_DIR))
	$Q $(PANDOC) --standalone --pdf-engine=xelatex --toc --top-level-division=chapter -o $(BIN_DIR)/${PROJECT}.pdf README.yaml README.md

# Start recipes
.PHONY: __start__
__start__: stop $(BIN_DIR)/$(GOOS)/$(PROJECT) | $(TMP_DIR) $(LOG_DIR); $(info $(M) Starting $(PROJECT) on $(GOOS))
	$(info $(M)   Check the logs in $(LOG_DIR) with `make logview`)
	$Q DEBUG=1 LOG_DESTINATION="$(LOG_DIR)/$(PROJECT).log" $(BIN_DIR)/$(GOOS)/$(PROJECT) & echo $$! > $(TMP_DIR)/$(PROJECT).pid

# publish recipes
.PHONY: __publish_init__ __publish_binaries__
__publish_init__:;
__publish_binaries__: archive
	$(info $(M) Uploading the binary packages...)
	$Q $(foreach archive, $(wildcard $(BIN_DIR)/*.tar.gz), go run . artifact upload $(archive) ;)
	$Q $(foreach archive, $(wildcard $(BIN_DIR)/*.zip), go run . artifact upload $(archive) ;)
	$(info $(M) Uploading the Debian packages...)
	$Q $(foreach archive, $(wildcard $(BIN_DIR)/*.deb), go run . artifact upload $(archive) ;)
	$(info $(M) Uploading the RPM packages...)
	$Q $(foreach archive, $(wildcard $(BIN_DIR)/*.rpm), go run . artifact upload $(archive) ;)

# archive recipes
.PHONY: __archive_init__ __archive_all__ __archive_chocolatey__ __archive_debian__ __archive_rpm__ __archive_snap__
__archive_init__:;      $(info $(M) Archiving binaries for application $(PROJECT))
__archive_all__: \
	$(BIN_DIR)/$(PACKAGE)_$(VERSION)_darwin_amd64.tar.gz \
	$(BIN_DIR)/$(PACKAGE)_$(VERSION)_darwin_arm64.tar.gz \
	$(BIN_DIR)/$(PACKAGE)_$(VERSION)_linux_amd64.tar.gz \
	$(BIN_DIR)/$(PACKAGE)_$(VERSION)_linux_arm64.tar.gz \
	$(BIN_DIR)/$(PACKAGE)-$(VERSION)-windows-amd64.zip \
	$(BIN_DIR)/$(PACKAGE)-$(VERSION)-windows-arm64.zip \
	;
__archive_chocolatey__: \
	packaging/chocolatey/tools/$(PACKAGE)-$(VERSION)-windows-amd64.7z \
	packaging/chocolatey/tools/$(PACKAGE)-$(VERSION)-windows-arm64.7z \
	;
__archive_debian__: \
	$(BIN_DIR)/$(PACKAGE)_$(VERSION)-$(REVISION)_amd64.deb \
	$(BIN_DIR)/$(PACKAGE)_$(VERSION)-$(REVISION)_arm64.deb \
	;
__archive_rpm__: \
	$(BIN_DIR)/$(PACKAGE)-$(VERSION)-$(REVISION).x86_64.rpm \
	$(BIN_DIR)/$(PACKAGE)-$(VERSION)-$(REVISION).aarch64.rpm \
	;

__archive_snap__: \
	$(BIN_DIR)/$(PACKAGE)_$(VERSION)_amd64.snap \
	;

$(BIN_DIR)/$(PACKAGE)_$(VERSION)_darwin_amd64.tar.gz: $(BIN_DIR)/darwin-amd64/$(PROJECT)
	$Q $(TAR) czf $@ -C $(<D) $(<F)
$(BIN_DIR)/$(PACKAGE)_$(VERSION)_darwin_arm64.tar.gz: $(BIN_DIR)/darwin-arm64/$(PROJECT)
	$Q $(TAR) czf $@ -C $(<D) $(<F)
$(BIN_DIR)/$(PACKAGE)_$(VERSION)_linux_amd64.tar.gz: $(BIN_DIR)/linux-amd64/$(PROJECT)
	$Q $(TAR) czf $@ -C $(<D) $(<F)
$(BIN_DIR)/$(PACKAGE)_$(VERSION)_linux_arm64.tar.gz: $(BIN_DIR)/linux-arm64/$(PROJECT)
	$Q $(TAR) czf $@ -C $(<D) $(<F)
$(BIN_DIR)/$(PACKAGE)-$(VERSION)-windows-amd64.zip: $(BIN_DIR)/windows-amd64/$(PROJECT).exe
	$Q $(ZIP) -9 -q --junk-paths $@ $<
$(BIN_DIR)/$(PACKAGE)-$(VERSION)-windows-arm64.zip: $(BIN_DIR)/windows-arm64/$(PROJECT).exe
	$Q $(ZIP) -9 -q --junk-paths $@ $<

packaging/chocolatey/tools/$(PACKAGE)-$(VERSION)-windows-amd64.7z: $(BIN_DIR)/windows-amd64/$(PROJECT).exe
	$Q $(7ZIP) a -r $@ $<
packaging/chocolatey/tools/$(PACKAGE)-$(VERSION)-windows-arm64.7z: $(BIN_DIR)/windows-arm64/$(PROJECT).exe
	$Q $(7ZIP) a -r $@ $<

$(BIN_DIR)/$(PACKAGE)_$(VERSION)-$(REVISION)_amd64.deb: packaging/nfpm.yaml $(BIN_DIR)/linux-amd64/$(PROJECT)
	$Q PLATFORM=linux-amd64 $(GOMPLATE) --file packaging/nfpm.yaml | $(NFPM) package --config - --target $(@D) --packager deb
$(BIN_DIR)/$(PACKAGE)_$(VERSION)-$(REVISION)_arm64.deb: packaging/nfpm.yaml $(BIN_DIR)/linux-arm64/$(PROJECT)
	$Q PLATFORM=linux-arm64 $(GOMPLATE) --file packaging/nfpm.yaml | $(NFPM) package --config - --target $(@D) --packager deb

$(BIN_DIR)/$(PACKAGE)-$(VERSION)-$(REVISION).x86_64.rpm: packaging/nfpm.yaml $(BIN_DIR)/linux-amd64/$(PROJECT)
	$Q PLATFORM=linux-amd64 $(GOMPLATE) --file packaging/nfpm.yaml | $(NFPM) package --config - --target $(@D) --packager rpm
$(BIN_DIR)/$(PACKAGE)-$(VERSION)-$(REVISION).aarch64.rpm: packaging/nfpm.yaml $(BIN_DIR)/linux-arm64/$(PROJECT)
	$Q PLATFORM=linux-arm64 $(GOMPLATE) --file packaging/nfpm.yaml | $(NFPM) package --config - --target $(@D) --packager rpm

$(BIN_DIR)/$(PACKAGE)_$(VERSION)_amd64.snap: packaging/snap/snapcraft.yaml
	$Q $(RM) $@
	$Q (cd packaging && snapcraft)
	$Q $(MOVE) packaging/$(@F) $(@D)

# build recipes for various platforms
.PHONY: __build_all__ __build_init__ __fetch_modules__
__build_init__:;     $(info $(M) Building application $(PROJECT))
__build_all__:       $(foreach platform, $(PLATFORMS), $(BIN_DIR)/$(platform)/$(PROJECT));
__fetch_modules__: ; $(info $(M) Fetching Modules...)
	$Q $(GO) mod download

$(BIN_DIR)/darwin-amd64: $(BIN_DIR) ; $(MKDIR)
$(BIN_DIR)/darwin-amd64/$(PROJECT): export GOOS=darwin
$(BIN_DIR)/darwin-amd64/$(PROJECT): export GOARCH=amd64
$(BIN_DIR)/darwin-amd64/$(PROJECT): $(GOFILES) $(ASSETS) | $(BIN_DIR)/darwin-amd64; $(info $(M) building application for darwin Intel)
	$Q $(GO) build $(if $V,-v) $(LDFLAGS) -o $@ .

$(BIN_DIR)/darwin-arm64: $(BIN_DIR) ; $(MKDIR)
$(BIN_DIR)/darwin-arm64/$(PROJECT): export GOOS=darwin
$(BIN_DIR)/darwin-arm64/$(PROJECT): export GOARCH=arm64
$(BIN_DIR)/darwin-arm64/$(PROJECT): $(GOFILES) $(ASSETS) | $(BIN_DIR)/darwin-arm64; $(info $(M) building application for darwin M1)
	$Q $(GO) build $(if $V,-v) $(LDFLAGS) -o $@ .

$(BIN_DIR)/linux-amd64: $(BIN_DIR) ; $(MKDIR)
$(BIN_DIR)/linux-amd64/$(PROJECT): export GOOS=linux
$(BIN_DIR)/linux-amd64/$(PROJECT): export GOARCH=amd64
$(BIN_DIR)/linux-amd64/$(PROJECT): $(GOFILES) $(ASSETS) | $(BIN_DIR)/linux-amd64; $(info $(M) building application for linux amd64)
	$Q $(GO) build $(if $V,-v) $(LDFLAGS) -o $@ .

$(BIN_DIR)/linux-arm64: $(BIN_DIR) ; $(MKDIR)
$(BIN_DIR)/linux-arm64/$(PROJECT): export GOOS=linux
$(BIN_DIR)/linux-arm64/$(PROJECT): export GOARCH=arm64
$(BIN_DIR)/linux-arm64/$(PROJECT): $(GOFILES) $(ASSETS) | $(BIN_DIR)/linux-arm64; $(info $(M) building application for linux arm64)
	$Q $(GO) build $(if $V,-v) $(LDFLAGS) -o $@ .

$(BIN_DIR)/windows-amd64: $(BIN_DIR) ; $(MKDIR)
$(BIN_DIR)/windows-amd64/$(PROJECT): $(BIN_DIR)/windows-amd64/$(PROJECT).exe;
$(BIN_DIR)/windows-amd64/$(PROJECT).exe: export GOOS=windows
$(BIN_DIR)/windows-amd64/$(PROJECT).exe: export GOARCH=amd64
$(BIN_DIR)/windows-amd64/$(PROJECT).exe: $(GOFILES) $(ASSETS) | $(BIN_DIR)/windows-amd64; $(info $(M) building application for windows amd64)
	$Q $(GO) build $(if $V,-v) $(LDFLAGS) -o $@ .

$(BIN_DIR)/windows-arm64: $(BIN_DIR) ; $(MKDIR)
$(BIN_DIR)/windows-arm64/$(PROJECT): $(BIN_DIR)/windows-arm64/$(PROJECT).exe;
$(BIN_DIR)/windows-arm64/$(PROJECT).exe: export GOOS=windows
$(BIN_DIR)/windows-arm64/$(PROJECT).exe: export GOARCH=arm64
$(BIN_DIR)/windows-arm64/$(PROJECT).exe: $(GOFILES) $(ASSETS) | $(BIN_DIR)/windows-arm64; $(info $(M) building application for windows arm64)
	$Q $(GO) build $(if $V,-v) $(LDFLAGS) -o $@ .

$(BIN_DIR)/pi:   $(BIN_DIR) ; $(MKDIR)
$(BIN_DIR)/pi/$(PROJECT): export GOOS=linux
$(BIN_DIR)/pi/$(PROJECT): export GOARCH=arm
$(BIN_DIR)/pi/$(PROJECT): export GOARM=6
$(BIN_DIR)/pi/$(PROJECT): $(GOFILES) $(ASSETS) | $(BIN_DIR)/pi; $(info $(M) building application for pi)
	$Q $(GO) build $(if $V,-v) $(LDFLAGS) -o $@ .

# Watch recipes
watch: watch-tools | $(TMP_DIR); @ ## Run a command continuously: make watch run="go test"
	@#$Q LOG=* $(YOLO) -i '*.go' -e vendor -e $(BIN_DIR) -e $(LOG_DIR) -e $(TMP_DIR) -c "$(run)"
	$Q nodemon \
	  --verbose \
	  --delay 5 \
	  --watch . \
	  --ext go \
	  --ignore .git/ --ignore bin/ --ignore log/ --ignore tmp/ \
	  --ignore './*.log' --ignore '*.md' \
	  --ignore go.mod --ignore go.sum  \
	  --exec "$(run) || exit 1"

# Download recipes
.PHONY: watch-tools coverage-tools
$(BIN_DIR)/yolo:      PACKAGE=github.com/azer/yolo
$(BIN_DIR)/gocov:     PACKAGE=github.com/axw/gocov/...
$(BIN_DIR)/gocov-xml: PACKAGE=github.com/AlekSi/gocov-xml
$(BIN_DIR)/nfpm:      PACKAGE=github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
$(BIN_DIR)/gomplate:  PACKAGE=github.com/hairyhenderson/gomplate/v4/cmd/gomplate@latest

watch-tools:    | $(YOLO)
coverage-tools: | $(GOCOV) $(GOCOVXML)

$(BIN_DIR)/%: | $(BIN_DIR) ; $(info $(M) installing $(PACKAGE)...)
	$Q tmp=$$(mktemp -d) ; \
	  env GOPATH=$$tmp GOBIN=$(BIN_DIR) $(GO) get $(PACKAGE) || status=$$? ; \
	  chmod -R u+w $$tmp ; rm -rf $$tmp ; \
	  exit $$status
