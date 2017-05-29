PKGS=$(shell go list ./... | grep -v /vendor)

compile:
	go install .

build:
	dep ensure

linux:
	docker run -i -t -v ${GOPATH}/src/github.com/kris-nova/klone:/go/src/github.com/kris-nova/klone -w /go/src/github.com/kris-nova/klone \
	-e TEST_KLONE_GITHUBTOKEN=${TEST_KLONE_GITHUBTOKEN} \
	-e TEST_KLONE_GITHUBUSER=${TEST_KLONE_GITHUBUSER} \
	-e TEST_KLONE_GITHUBPASS=${TEST_KLONE_GITHUBPASS} \
    --rm golang:1.8.1

test:
	docker run -v ${GOPATH}/src/github.com/kris-nova/klone:/go/src/github.com/kris-nova/klone -w /go/src/github.com/kris-nova/klone \
	-e TEST_KLONE_GITHUBTOKEN=${TEST_KLONE_GITHUBTOKEN} \
	-e TEST_KLONE_GITHUBUSER=${TEST_KLONE_GITHUBUSER} \
	-e TEST_KLONE_GITHUBPASS=${TEST_KLONE_GITHUBPASS} \
	--rm golang:1.8.1 make local-test

local-test:
	@go test $(PKGS)

test-race:
	@go test -race $(PKGS)