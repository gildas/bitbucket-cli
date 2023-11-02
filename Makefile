-include .env

# Goodies
V = 0
Q = $(if $(filter 1,$V),,@)
E := 
S := $E $E
M = $(shell printf "\033[34;1m▶\033[0m")
rwildcard = $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))

# Folders
BIN_DIR ?= $(CURDIR)/bin
LOG_DIR ?= log
TMP_DIR ?= tmp
COV_DIR ?= tmp/coverage

# Version, branch, and project
BRANCH    != git symbolic-ref --short HEAD
COMMIT    != git rev-parse --short HEAD
STAMP     != date +%Y%m%d%H%M%S
BUILD     := "$(STAMP).$(COMMIT)"
VERSION   != awk '/^var +VERSION +=/{gsub("\"", "", $$4) ; print $$4}' version.go
ifeq ($VERSION,)
VERSION   != git describe --tags --always --dirty="-dev"
endif
PROJECT   != awk '/^const +APP += +/{gsub("\"", "", $$4); print $$4}' version.go
ifeq (${PROJECT},)
PROJECT   != basename "$(PWD)"
endif
PLATFORMS ?= darwin linux windows pi

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
DOCKER  ?= docker
PANDOC  ?= pandoc

# Flags
#MAKEFLAGS += --silent
# GO
LDFLAGS = -ldflags "-X main.commit=$(COMMIT) -X main.branch=$(BRANCH) -X main.stamp=$(STAMP)"
ifneq ($(what),)
TEST_ARG := -run '$(what)'
else
TEST_ARG :=
endif

# Docker
DOCKER_FILE         ?= Dockerfile
ifneq ("$(wildcard chart/values.yaml)", "")
DOCKER_REGISTRY     != awk '/^  registry:/{print $$2}'   chart/values.yaml
DOCKER_REPOSITORY   != awk '/^  repository:/{print $$2}' chart/values.yaml
else
DOCKER_REGISTRY     ?= gildas
DOCKER_REPOSITORY   ?= cantina
endif
DOCKER_IMAGE         = $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY)
ifneq ($(BRANCH), master)
ifneq ("$(wildcard Dockerfile.$(BRANCH))", "")
DOCKER_FILE              ?= Dockerfile.$(BRANCH)
else ifneq ("$(wildcard Dockerfile.dev)", "")
DOCKER_FILE              ?= Dockerfile.dev
endif
DOCKER_IMAGE_VERSION     := $(VERSION)-$(STAMP)-$(COMMIT)
DOCKER_IMAGE_TAG         := dev
else
DOCKER_IMAGE_VERSION     := $(VERSION)
DOCKER_IMAGE_VERSION_MIN := $(subst $S,.,$(wordlist 1,2,$(subst .,$S,$(DOCKER_IMAGE_VERSION)))) # Removes the last .[0-9]+ part of the version
DOCKER_IMAGE_VERSION_MAJ := $(subst $S,.,$(wordlist 1,1,$(subst .,$S,$(DOCKER_IMAGE_VERSION)))) # Removes the 2 last .[0-9]+ parts of the version
DOCKER_IMAGE_TAG         := latest
endif
# Check the GOPROXY for localhost, so the docker rule can call the Athens container on the localhost
ifneq ($(findstring localhost, $(GOPROXY)),)
ifeq ($(GOOS), linux)
DOCKER_FLAGS += --network=host
DOCKER_GOPROXY := $(subst localhost,127.0.0.1,$(GOPROXY))
else
  # Docker for Windows, Docker for Mac
DOCKER_GOPROXY := $(subst localhost, host.docker.internal,$(GOPROXY))
endif
endif

# Kubernetes
#K8S_NAMESPACE =
K8S_APP = $(subst _,-,$(PROJECT))
ifneq ($(K8S_NAMESPACE),)
K8S_FLAGS += --namespace $(K8S_NAMESPACE)
endif

# Main Recipes
.PHONY: all build dep deploy docker fmt gendoc help lint logview publish run start stop test version vet watch

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

ifneq ("$(wildcard chart/values.yaml)", "")
deploy: __k8s_deploy_init__ __k8s_deploy__ __k8s_refresh_pods__; @ ## Deploy to a Kubernetes Cluster
endif

publish: docker __publish_init__; @ ## Publish the Docker image to the Repository
ifeq ($(DOCKER_REGISTRY),)
	$(error     DOCKER_REGISTRY is undefined, we cannot push the Docker image $(DOCKER_REPOSITORY))
else
	$Q $(DOCKER) push $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION)
ifneq ($(DOCKER_IMAGE_VERSION_MIN),)
	$Q $(DOCKER) push $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION_MIN)
endif
ifneq ($(DOCKER_IMAGE_VERSION_MAJ),)
	$Q $(DOCKER) push $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION_MAJ)
endif
	$Q $(DOCKER) push $(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG)
endif

docker: $(TMP_DIR)/__docker_$(BRANCH)__; @ ## Build the Docker image

build: __build_init__ __build_all__; @ ## Build the application for all platforms

dep:; $(info $(M) Updating Modules...) @ ## Updates the GO Modules
	$Q $(GO) get -u -t ./...
	$Q $(GO) mod tidy

lint:;  $(info $(M) Linting application...) @ ## Lint Golang files
	$Q $(GOLINT) run *.go

fmt:; $(info $(M) Formatting the code...) @ ## Format the code following the go-fmt rules
	$Q $(GO) fmt -n *.go

vet:; $(info $(M) Vetting application...) @ ## Run go vet
	$Q $(GO) vet *.go

run:; $(info $(M) Running application...) @  ## Execute the application
	$Q $(GO) run . | $(LOGGER)

logview:; @ ## Open the project log and follows it
	$Q tail -F $(LOG_DIR)/$(PROJECT).log | $(LOGGER)

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
$(BIN_DIR): ; @mkdir -p $@
$(TMP_DIR): ; @mkdir -p $@
$(LOG_DIR): ; @mkdir -p $@
$(COV_DIR): ; @mkdir -p $@

# Documentation recipes
__gendoc_init__:; $(info $(M) Building the documentation...)

$(BIN_DIR)/$(PROJECT).pdf: README.md ; $(info $(M) Generating PDF documentation in $(BIN_DIR))
	$Q $(PANDOC) --standalone --pdf-engine=xelatex --toc --top-level-division=chapter -o $(BIN_DIR)/${PROJECT}.pdf README.yaml README.md

# Start recipes
.PHONY: __start__
__start__: stop $(BIN_DIR)/$(GOOS)/$(PROJECT) | $(TMP_DIR) $(LOG_DIR); $(info $(M) Starting $(PROJECT) on $(GOOS))
	$(info $(M)   Check the logs in $(LOG_DIR) with `make logview`)
	$Q DEBUG=1 LOG_DESTINATION="$(LOG_DIR)/$(PROJECT).log" $(BIN_DIR)/$(GOOS)/$(PROJECT) & echo $$! > $(TMP_DIR)/$(PROJECT).pid

# Docker recipes
# @see https://www.gnu.org/software/make/manual/html_node/Empty-Targets.html
# TODO: Pass LDFLAGS!!!
ifeq ($(BRANCH), master)
$(TMP_DIR)/__docker_$(BRANCH)__: $(GOFILES) $(ASSETS) $(DOCKER_FILE) | $(TMP_DIR); $(info $(M) Building the Docker Image...)
	$(info $(M)  Image: $(DOCKER_IMAGE), Version: $(DOCKER_IMAGE_VERSION), Tag: $(DOCKER_IMAGE_TAG), Branch: $(BRANCH))
	$Q DOCKER_BUILDKIT=1 $(DOCKER) build \
		$(DOCKER_FLAGS) \
		--build-arg GOPROXY=$(DOCKER_GOPROXY) \
		--label "org.opencontainers.image.version"="$(DOCKER_IMAGE_VERSION)" \
		--label "org.opencontainers.image.revision"="$(COMMIT)" \
		--label "org.opencontainers.image.created"="$(NOW)" \
		-t $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION) .
	$Q $(DOCKER) tag $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION) $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION_MIN)
	$Q $(DOCKER) tag $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION) $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION_MAJ)
	$Q $(DOCKER) tag $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION) $(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG)
	$Q touch $@
else
$(TMP_DIR)/__docker_$(BRANCH)__: $(GOFILES) $(ASSETS) $(DOCKER_FILE) build | $(TMP_DIR); $(info $(M) Building the Docker Image...)
	$(info $(M)  Image: $(DOCKER_IMAGE), Version: $(DOCKER_IMAGE_VERSION), Tag: $(DOCKER_IMAGE_TAG), Branch: $(BRANCH))
	$Q DOCKER_BUILDKIT=1 $(DOCKER) build \
		$(DOCKER_FLAGS) \
		-f $(DOCKER_FILE) \
		-f "$(DOCKER_FILE)" \
		--label "org.opencontainers.image.version"="$(DOCKER_IMAGE_VERSION)" \
		--label "org.opencontainers.image.revision"="$(COMMIT)" \
		--label "org.opencontainers.image.created"="$(NOW)" \
		-t $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION) .
	$Q $(DOCKER) tag $(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION) $(DOCKER_IMAGE):$(DOCKER_IMAGE_TAG)
	$Q touch $@
endif

__publish_init__:; $(info $(M) Pushing the Docker Image $(DOCKER_IMAGE) to $(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY)...)

.PHONY: __docker_save__
__docker_save__: $(TMP_DIR)/image.$(BRANCH).$(COMMIT).tar

$(TMP_DIR)/image.$(BRANCH).$(COMMIT).tar: $(TMP_DIR)/__docker_$(BRANCH)__ | $(TMP_DIR); $(info $(M)   Saving Docker image)
	$Q $(DOCKER) save $(DOCKER_IMAGE) > $(TMP_DIR)/image.$(BRANCH).$(COMMIT).tar

# Kubernetes recipes
.PHONY: __k8s_deploy_init__ __k8s_deploy__ __microk8s_import__ __remote_docker__
__k8s_deploy_init__:;  $(info $(M) Deploying to Kubernetes)

ifeq ($(K8S_TYPE), microk8s)
ifeq ($(K8S_REMOTE), 1)
__k8s_deploy__:; $(info (M) Error: Not implemented)
else
__k8s_deploy__: __docker_save__; $(info $(M)   Importing Image to MicroK8s)
	$Q microk8s.ctr image import $(TMP_DIR)/image.$(BRANCH).$(COMMIT).tar
endif
else ifeq ($(K8S_TYPE), docker)
ifeq ($(K8S_REMOTE), 1)
__k8s_deploy__: $(TMP_DIR)/__docker_$(BRANCH)__ | $(TMP_DIR); $(info $(M)    Sending Docker image to $(DOCKER_HOST))
	$Q $(DOCKER) save $(DOCKER_IMAGE) | ssh $(DOCKER_USER)@$(DOCKER_HOST) docker load
else
__k8s_deploy__:;
endif
else
__k8s_deploy__:; $(info $(M) K8S import is not set up) # Let's pretend there is nothing to do
endif

__k8s_refresh_pods__:; $(info $(M)   Restarting Pods for application $(K8S_APP))
	$Q $(KUBECTL) delete pod $(K8S_FLAGS) --selector app.kubernetes.io/name=$(K8S_APP)
	$Q $(KUBECTL) rollout status $(K8S_FLAGS) \
		`$(KUBECTL) get deployments $(K8S_FLAGS) --selector app.kubernetes.io/name=$(K8S_APP) --output name`

# build recipes for various platforms
.PHONY: __build_all__ __build_init__ __fetch_modules__
__build_init__:;     $(info $(M) Building application $(PROJECT))
__build_all__:       __fetch_modules__ $(foreach platform, $(PLATFORMS), $(BIN_DIR)/$(platform)/$(PROJECT));
__fetch_modules__: ; $(info $(M) Fetching Modules...)
	$Q $(GO) mod download

$(BIN_DIR)/darwin: $(BIN_DIR) ; @mkdir -p $@
$(BIN_DIR)/darwin/$(PROJECT): $(GOFILES) $(ASSETS) | $(BIN_DIR)/darwin; $(info $(M) building application for darwin)
	$Q CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(if $V,-v) $(LDFLAGS) -o $@ .

$(BIN_DIR)/linux:   $(BIN_DIR) ; @mkdir -p $@
$(BIN_DIR)/linux/$(PROJECT): $(GOFILES) $(ASSETS) | $(BIN_DIR)/linux; $(info $(M) building application for linux)
	$Q CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(if $V,-v) $(LDFLAGS) -o $@ .

$(BIN_DIR)/windows: $(BIN_DIR) ; @mkdir -p $@
$(BIN_DIR)/windows/$(PROJECT): $(BIN_DIR)/windows/$(PROJECT).exe;
$(BIN_DIR)/windows/$(PROJECT).exe: $(GOFILES) $(ASSETS) | $(BIN_DIR)/windows; $(info $(M) building application for windows)
	$Q CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(if $V,-v) $(LDFLAGS) -o $@ .

$(BIN_DIR)/pi:   $(BIN_DIR) ; @mkdir -p $@
$(BIN_DIR)/pi/$(PROJECT): $(GOFILES) $(ASSETS) | $(BIN_DIR)/pi; $(info $(M) building application for pi)
	$Q CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 $(GO) build $(if $V,-v) $(LDFLAGS) -o $@ .

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

watch-tools:    | $(YOLO)
coverage-tools: | $(GOCOV) $(GOCOVXML)

$(BIN_DIR)/%: | $(BIN_DIR) ; $(info $(M) installing $(PACKAGE)...)
	$Q tmp=$$(mktemp -d) ; \
	  env GOPATH=$$tmp GOBIN=$(BIN_DIR) $(GO) get $(PACKAGE) || status=$$? ; \
	  chmod -R u+w $$tmp ; rm -rf $$tmp ; \
	  exit $$status
