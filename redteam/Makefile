BINARY ?= bin/redteam
MOCKBIN ?= bin/mocktarget

.PHONY: build build-mock test run-mock run-list clean

build:
	go build -ldflags "-X github.com/llm-redteam/redteam/internal/version.Version=dev" -o $(BINARY) ./cmd/redteam

build-mock:
	go build -o $(MOCKBIN) ./cmd/mocktarget

test:
	go test ./...

run-list: build
	./$(BINARY) list

run-mock: build
	./$(BINARY) run --config examples/redteam-mock.yaml --output report.json --format json --mock --fail-on-findings

run-html: build
	./$(BINARY) run --config examples/redteam-mock.yaml --output report.html --format html --mock

clean:
	rm -rf bin/ report.json report.html

# Terminal 1: make run-mockserver-vulnerable
# Terminal 2: make run-against-mockserver
run-mockserver-vulnerable: build-mock
	./$(MOCKBIN) -addr 127.0.0.1:8765 -persona vulnerable

run-against-mockserver: build
	./$(BINARY) run --config examples/redteam-local-mockserver.yaml --output report.html --format html --fail-on-findings
