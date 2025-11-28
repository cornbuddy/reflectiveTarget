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

.PHONY: test
test:
	ENVIRONMENT=development \
	GOTOOLCHAIN=go1.25.3+auto \
		go test -coverprofile=$(COVERPROFILE) $(PACKAGES)
	go tool cover -html $(COVERPROFILE) -o $(COVERREPORT)
	- xdg-open $(COVERREPORT)

.PHONY: deps
deps: update tidy

.PHONY: update
update:
	go get -u ./...
