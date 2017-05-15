PKGS=$(shell go list ./... | grep -v /vendor)

compile:
	go install .

build:
	dep ensure

test:
	@go test $(PKGS)

test-race:
	@go test -race $(PKGS)