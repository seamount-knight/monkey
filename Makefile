OS = Linux

VERSION = 0.0.1

CURDIR = $(shell pwd)
SOURCEDIR = $(CURDIR)
COVER = $($3)

ECHO = echo
RM = rm -rf
MKDIR = mkdir

# If the first argument is "cover"...
ifeq (cover,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif
ifeq (mysql,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif
ifeq (psql,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

.PHONY: test

default: test

test:
	go test -cover=true $(PACKAGES)

race:
	go test -cover=true -race $(PACKAGES)

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt $(PACKAGES)

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint .

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet $(PACKAGES)

setup:
	@$(ECHO) "Installing gomock..."
	go get -u github.com/golang/mock/gomock
	go get -u github.com/golang/mock/mockgen
	@$(ECHO) "Installing govendor..."
	go get -u github.com/kardianos/govendor
	@$(ECHO) "Installing cobra..."
	go get -v github.com/spf13/cobra/cobra

add:
	govendor add +external

all: test

PACKAGES = $(shell go list ./... | grep -v ./vendor/)
BUILD_PATH = $(shell if [ "$(ALAUDACI_DEST_DIR)" != "" ]; then echo "$(ALAUDACI_DEST_DIR)" ; else echo "$(PWD)"; fi)

cover: collect-cover-data test-cover-html open-cover-html

collect-cover-data:
	echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES),\
		go test -v -coverprofile=coverage.out -covermode=count $(pkg);\
		if [ -f coverage.out ]; then\
			tail -n +2 coverage.out >> coverage-all.out;\
		fi;)

test-cover-html:
	go tool cover -html=coverage-all.out -o coverage.html

test-cover-func:
	go tool cover -func=coverage-all.out

open-cover-html:
	open coverage.html

build:
	@$(ECHO) "Will build on "$(BUILD_PATH)
	go build -ldflags "-w -s" -v -o $(BUILD_PATH)/bin/monkey monkey

clean:
	./scripts/clean.sh

compile: test build

psql:
	docker-compose -p gohellop -f devbox/psql-compose.yaml $(RUN_ARGS)

docker-build:
	docker build -t gohello .

mysql:
	docker-compose -p gohellom -f devbox/mysql-compose.yaml $(RUN_ARGS)

help:
	@$(ECHO) "Targets:"
	@$(ECHO) "all				- test"
	@$(ECHO) "setup				- install necessary libraries"
	@$(ECHO) "test				- run all unit tests"
	@$(ECHO) "cover [package]			- generates and opens unit test coverage report for a package"
	@$(ECHO) "race				- run all unit tests in race condition"
	@$(ECHO) "add				- runs govendor add +external command"
	@$(ECHO) "build				- build and exports using ALAUDACI_DEST_DIR"
	@$(ECHO) "clean				- remove test reports and compiled package from this folder"
	@$(ECHO) "compile				- test and build - one command for CI"
	@$(ECHO) "psql				- runs a psql compose command. e.g: make psql up"
	@$(ECHO) "mysql				- runs a mysql compose command e.g: make mysql down"
	@$(ECHO) "docker-build			- builds an image with this folder's Dockerfile"


