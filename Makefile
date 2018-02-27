APP ?= rumpel

bin := $(GOPATH)/bin
src := ./cmd/$(APP)/main.go
dst := $(bin)/$(APP)

gometalinter := $(bin)/gometalinter

.PHONY: deps test lint build

$(dst): test
	@echo "===> Building app..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -installsuffix nocgo -o $(dst) $(src)

test: lint
	@echo "===> Testing packages..."
	@echo "mode: count" > coverage-all.out
	@for pkg in $(shell go list ./... | grep -v /vendor/); do \
		go test -coverprofile=coverage.out -covermode=count $$pkg -timeout 30ms -check.v && \
		tail -n +2 coverage.out >> ./coverage-all.out; \
	done

lint: $(gometalinter)
	@echo "===> Executing linter..."
	@gometalinter ./... --tests --vendor

$(gometalinter):
	@echo "===> Installing gometalinter..."
	go get -u -v github.com/alecthomas/gometalinter
	gometalinter -i

cover:
	@echo "===> Executing cover...."
	@go tool cover -html=coverage-all.out

