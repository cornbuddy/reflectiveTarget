.PHONY: all
all: lint test

.PHONY: lint
lint: fmt vet tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy

COVERPROFILE = cover.out
COVERREPORT = cover.html
PACKAGES = ./...
TESTS = ^.+$
GO_TEST_CMD = go test -vet=all -count=1 \
-run $(TESTS) \
-coverprofile=$(COVERPROFILE) \
$(PACKAGES)

.PHONY: test
test:
	ENVIRONMENT=development \
	GOTOOLCHAIN=$(shell go env GOVERSION)+auto \
		$(GO_TEST_CMD)
	go tool cover -html $(COVERPROFILE) -o $(COVERREPORT)
	- xdg-open $(COVERREPORT)

.PHONY: deps
deps: update tidy

.PHONY: update
update:
	go get -u ./...
