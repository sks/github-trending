GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

.DEFAULT_GOAL = all

BINARY ?= rssfeed
ARCH ?= amd64

.PHONY: all

GO_PROJECT = github.com/sks/github-trending
BUILD_DEST = build
SERVERLESS_ZIP = ${BUILD_DEST}/serverless.zip
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
VERSION = ${GIT_COMMIT}-${GIT_DIRTY}


LDFLAGS += -w -s

DOCKER_BUILD_ARGS=--build-arg VERSION=${GIT_COMMIT}-${GIT_DIRTY} \
	--build-arg GOARCH=${ARCH}

GOFLAGS := -v -ldflags "$(LDFLAGS)"

## Download dependencies and the run unit test and build the binary
all: test all/amd64 dockerize

all/amd64: setup test clean build

## Run the application on local machine
run: start

start:
	GO111MODULE=on go run cmd/${BINARY}/main.go

## Clean the dist directory
clean:
	@echo "Cleaning the $(BUILD_DEST)"
	@rm -rf $(BUILD_DEST)

## Update depdendencies
dep:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

## Run the ginkgo test
unit:
	go test ./...

fakes:
	@find . -iname interfaces.go -exec  go generate {}  \;

## Run unit tests : Defaults to AMD64
test: fakes unit

## download dependencies to run this project
setup:
	@which counterfeiter > /dev/null || GO111MODULE=off go get -u github.com/maxbrunsfeld/counterfeiter

BADFMT=$(shell find . -name '*.go' ! -path '*/vendor/*' -exec gofmt -l {} \+)

ifeq ($(BADFMT),)
format: $(BADFMT)
	@echo "All files are formatted properly. Good job"
else
format: $(BADFMT)
	@echo "Some files are not formatted properly!!. Fix it Using the command below."
	@echo "\n\t${GREEN}find . -name '*.go' ! -path '*/vendor/*' -exec gofmt -w {} \+${RESET}\n"
	exit 1
endif

## Build the linux binary
build: clean format binaries serverless

serverless:
	@$(MAKE) GOOS=linux BINARY=serverless build-amd64
	cp cmd/serverless/serverless/helper.go ${BUILD_DEST}/
	zip -r -j ${SERVERLESS_ZIP} ${BUILD_DEST}/helper.go ${BUILD_DEST}/serverless_linux_amd64

deploy: serverless
	@which gcloud > /dev/null  || brew cask install google-cloud-sdk || echo "Not sure how to install google-cloud-sdk"
	@which gsutil > /dev/null  || brew cask install google-cloud-sdk || echo "Not sure how to install google-cloud-sdk"
	gsutil cp ${SERVERLESS_ZIP} gs://${STORAGE_LOCATION}/${VERSION}.zip
	gcloud functions deploy github-trending \
		--runtime go111 \
		--set-env-vars=VERSION=$$(date +%s) \
		--set-env-vars=COMMIT_SHA=${VERSION} \
		--set-env-vars=COMMAND_TO_RUN=serverless_linux_amd64 \
		--source=gs://${STORAGE_LOCATION}/${VERSION}.zip \
		--project=${PROJECT_ID}


binaries:
	@$(MAKE) BINARY=cli build-amd64
	@$(MAKE) BINARY=rssfeed build-amd64
	@$(MAKE) BINARY=serverless build-amd64

# ## Build the docker images
dockerize:
# 	@$(MAKE) dockerize-arch-amd64

# dockerize-arch-%:
# 	@$(MAKE) _dockerize ARCH=$* BINARY=cli

# dockerize-binary-%:
# 	@$(MAKE) _dockerize ARCH=amd64 BINARY=$*

build-%:
	@mkdir -p $(BUILD_DEST) > /dev/null
	GOARCH=$* \
		go build $(GOFLAGS) -o $(BUILD_DEST)/${BINARY}_${GOOS}_$* ./cmd/${BINARY}/

_dockerize:
	docker build $(DOCKER_BUILD_ARGS) \
		-t ${BINARY}:${ARCH}-latest \
		-f cmd/${BINARY}/Dockerfile .

## Prints the version info about the project
info:
	 @echo "Git Commit:        ${GIT_COMMIT}"
	 @echo "Git Tree State:    ${GIT_DIRTY}"
	 @echo "Version: ${VERSION}"

## Prints this help command
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET}: ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)