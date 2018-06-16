all: pre test

pre: fmt vet lint

fmt:
	go fmt ./...

vet:
	go vet $$(go list ./... | grep -v /vendor/)

lint:
	golint $$(go list ./... | grep -v /vendor/)

test:
	go test -v -race -cover ./...

dep-init:
	dep ensure

dep-update:
	dep ensure -update

.PHONY: all pre fmt vet lint test dep-init dep-update
