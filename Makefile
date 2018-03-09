APP ?= rumpel
REPORT_FILE ?= report

bin := $(GOPATH)/bin
src := ./cmd/$(APP)/main.go
dst := $(bin)/$(APP)

gometalinter := $(bin)/gometalinter
gojunitreport := $(bin)/go-junit-report

.PHONY: deps test lint build

$(dst): test
	@echo "===> Building app..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -installsuffix nocgo -o $(dst) $(src)

test: $(gojunitreport) lint
	@echo "===> Testing packages..."
	@2>&1 go test -v -cover $(shell go list ./... | grep -v /vendor/) | tee tests.out
	@cat tests.out | go-junit-report -set-exit-code=1 > $(REPORT_FILE).xml

lint: $(gometalinter)
	@echo "===> Executing linter..."
	@gometalinter ./... --tests --vendor

cover:
	@echo "===> Coveraging code from packages..."
	@echo "mode: count" > coverage-all.out
	@for pkg in $(shell go list ./... | grep -v /vendor/); do \
		go test -coverprofile=coverage.out -covermode=count $$pkg -timeout 30ms && \
		tail -n +2 coverage.out >> ./coverage-all.out; \
	done
	@go tool cover -html=coverage-all.out

$(gometalinter):
	@echo "===> Installing gometalinter..."
	go get -u -v github.com/alecthomas/gometalinter
	gometalinter -i

$(gojunitreport):
	@echo "===> Installing go-junit-report..."
	go get -u -v github.com/jstemmer/go-junit-report

