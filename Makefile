GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "*/vendor/*" 2>/dev/null)
GOPACKAGES = $(shell go list ./...)
GOIMPORTS_REPO = golang.org/x/tools/cmd/goimports
DIST = $(shell basename $(PWD))
REPO_URL = github.com/nistal97/${DIST}
BUILD_DIR = .build
ANTLR_PATH = ~/dev/tools/antlr-4.11.1-complete.jar
CLASSPATH = ${ANTLR_PATH}:hack/neo4j
export KO_DOCKER_REPO = 127.0.0.1:31676/${DIST}

T_GOOS = linux
T_GOARCH = amd64
GOOS = win32
GOARCH = amd64
ifeq ($(OS),Windows_NT)
	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		OSFLAG += -D IA32
	endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS = linux
	endif
	ifeq ($(UNAME_S),Darwin)
		GOOS = darwin
	endif
		UNAME_P := $(shell uname -p)
	ifneq ($(filter %86,$(UNAME_P)),)
		OSFLAG += -D IA32
	endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		OSFLAG += -D arm64
	endif
endif


.PHONY: clean
clean:
	@rm -rf ${BUILD_DIR}

.PHONY: init
init: ## Initialize workspace
	@echo "Initializing..."
	@mkdir -p cmd pkg
	@echo ${REPO_URL}
	@go mod init ${REPO_URL}

.PHONY: fmt
fmt: ## Ensures all go files are properly formatted.
	@echo "Formatting..."
	@GO111MODULE=off go get -u ${GOIMPORTS_REPO}
	@goimports -l -w ${GOFILES_NOVENDOR}

.PHONY: parse
parse:
	#@java -cp ${CLASSPATH} org.antlr.v4.Tool hack/neo4j/Cypher.g4
	#@javac -cp ${CLASSPATH} hack/neo4j/*.java
	#@java -cp ${CLASSPATH} org.antlr.v4.runtime.misc.TestRig $* Cypher OC_Match -gui
	@java -cp ${CLASSPATH} org.antlr.v4.Tool -Dlanguage=Go -o pkg/parser/calc hack/calc/Calc.g4
	@java -cp ${CLASSPATH} org.antlr.v4.Tool -Dlanguage=Go -o pkg/parser/neo4j hack/neo4j/Cypher.g4

.PHONY: get
get: export GO111MODULE=on
get: export GOPRIVATE=github.coupang.net
get:
	#@go get github.com/machinebox/graphql
    #@go get github.com/antlr/antlr4/runtime/Go/antlr/v4

.PHONY: vendor
vendor: export GO111MODULE=on
vendor: export GOPRIVATE=github.coupang.net
vendor: ## Ensures all go module dependencies are synced and copied to vendor
	@git config --global url."git@github.coupang.net:".insteadOf "https://github.coupang.net/"
	@echo "Updating module dependencies..."
	@go mod tidy
	@go mod vendor

.PHONY: build
build: vendor
	@echo "Building binary on ${OSFLAG}..."
	@mkdir -p ${BUILD_DIR}
	@GO111MODULE=on GOOS=${GOOS} GOARCH=${GOARCH} go build -v -mod vendor -o ${BUILD_DIR} ${GO_LDFLAGS} ./cmd

.PHONY: img
img: build
	@echo "Building image..."
	@mkdir -p ${BUILD_DIR}
	@ko build ./cmd --bare --sbom=none -t test --platform=${T_GOOS}/${T_GOARCH}

