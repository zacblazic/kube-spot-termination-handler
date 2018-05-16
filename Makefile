SHELL := /bin/bash

.PHONY: dependencies build install test

dependencies:
	dep ensure

build:
	@mkdir -p bin/
	go build -o ./bin/kube-spot-termination-handler

install:
	go install

test:
	go test -v ./...
