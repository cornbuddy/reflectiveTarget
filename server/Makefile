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

.PHONY: test
test: lint
	go test -cover ./...

.PHONY: deps
deps: update tidy

.PHONY: update
update:
	go get -u ./...
