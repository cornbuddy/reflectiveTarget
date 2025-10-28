.PHONY: run
run:
	go run main.go

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

COVERPROFILE=cover.out
COVERREPORT=cover.html

.PHONY: test
test: lint
	go test -coverprofile=$(COVERPROFILE) ./...
	go tool cover -html $(COVERPROFILE) -o $(COVERREPORT)
	- xdg-open $(COVERREPORT)

.PHONY: deps
deps: update tidy

.PHONY: update
update:
	go get -u ./...
