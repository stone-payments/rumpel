PROJECT := rumpel
ENTRYPOINT_DIR := ./cmd/$(PROJECT)
ENTRYPOINT_FILE := $(ENTRYPOINT_DIR)/main.go

BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter
GVT := $(BIN_DIR)/gvt

.PHONY: deps test lint build

$(BIN_DIR)/$(PROJECT): test
	CGO_ENABLED=0 go install $(ENTRYPOINT_DIR)

unit-test:
	@go test -cover -timeout 30ms $(shell go list ./... | grep -v /vendor/) -check.v

test: lint unit-test

lint: $(GOMETALINTER)
	@gometalinter ./... --tests --vendor

$(GOMETALINTER):
	go get -u -v github.com/alecthomas/gometalinter
	gometalinter -i

deps: $(GVT)
	gvt restore -connections 4

run:
	go run -race $(ENTRYPOINT_FILE)

$(GVT):
	go get -u -v github.com/FiloSottile/gvt

cover:
	@echo "mode: count" > coverage-all.out
	@for pkg in $(shell go list ./... | grep -v /vendor/); do \
		go test -coverprofile=coverage.out -covermode=count $$pkg -timeout 30ms -check.v && \
		tail -n +2 coverage.out >> ./coverage-all.out; \
	done
	@go tool cover -html=coverage-all.out

