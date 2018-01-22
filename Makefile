#! /usr/bin/make
#
# Makefile for blockchain-simulation
#
# Targets:
# - clean     delete all generated files and executable
# - dependencies    install or update dependencies
# - build     compile executable
# - run     compile executable
#
# Meta targets:
# - all is the default target, it runs clean, dependencies, build.

APP?=blockchain-simulation

all: clean dependencies build

dependencies:
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure

clean:
	@rm -f ${APP}

build:
	@go build -a -o ${APP}

