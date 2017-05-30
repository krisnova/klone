PKGS=$(shell go list ./... | grep -v /vendor)


compile:
	go install .

build:
	dep ensure


linux: shell

shell:
	docker run \
	-w /go/src/github.com/kris-nova/klone \
	-v ${GOPATH}/src/github.com/kris-nova/klone:/go/src/github.com/kris-nova/klone \
	-v ${HOME}/.ssh:/root/.ssh \
	-e TEST_KLONE_GITHUBTOKEN=${TEST_KLONE_GITHUBTOKEN} \
	-e TEST_KLONE_GITHUBUSER=${TEST_KLONE_GITHUBUSER} \
	-e TEST_KLONE_GITHUBPASS=${TEST_KLONE_GITHUBPASS} \
	-e KLONE_GITHUBTOKEN=${KLONE_GITHUBTOKEN} \
	-e KLONE_GITHUBUSER=${KLONE_GITHUBUSER} \
	-e KLONE_GITHUBPASS=${KLONE_GITHUBPASS} \
    --rm golang:1.8.1

test:
	@go test $(PKGS)