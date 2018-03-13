APP ?= rumpel
REPORT_FILE ?= report
COVERAGE_FILE ?= cover

bin := $(GOPATH)/bin
src := ./cmd/$(APP)/main.go
dst := $(bin)/$(APP)

gometalinter := $(bin)/gometalinter
gocover-cobertura := $(bin)/gocover-cobertura
go-junit-report := $(bin)/go-junit-report

.PHONY: deps test lint build

$(dst): cover
	@echo "===> Building app..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -installsuffix nocgo -o $(dst) $(src)

cover: $(gocover-cobertura) test
	@echo "===> Executing cover..."
	@echo "mode: count" > coverall.out
	@touch c.out
	@for pkg in $(shell go list ./... | grep -v /vendor/); do \
		go test -coverprofile=c.out -covermode=count $$pkg -timeout 30ms && \
		tail -n +2 c.out >> coverall.out; \
	done
	@gocover-cobertura < coverall.out > $(COVERAGE_FILE).xml
	@go tool cover -html=coverall.out -o cover.html

test: $(go-junit-report) lint
	@echo "===> Testing packages..."
	@go test -v -cover $(shell go list ./... | grep -v /vendor/) | go-junit-report -set-exit-code=1 > $(REPORT_FILE).xml

lint: $(gometalinter)
	@echo "===> Executing linter..."
	@gometalinter ./... --tests --vendor

$(gometalinter):
	@echo "===> Installing gometalinter..."
	go get -u -v github.com/alecthomas/gometalinter
	gometalinter -i

$(gocover-cobertura):
	@echo "===> Installing gocover-cobertura..."
	go get -u -v github.com/t-yuki/gocover-cobertura

$(go-junit-report):
	@echo "===> Installing go-junit-report..."
	go get -u -v github.com/jstemmer/go-junit-report

