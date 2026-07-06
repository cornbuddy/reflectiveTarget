include ../go.Makefile

.PHONY: generate
generate:
	go tool templ generate

test: generate

lint: generate

vet: generate

build: generate
	go build \
		-a -v \
		-gcflags=all="-l -B" \
		-ldflags="-w -s" \
		-o reflecitve-target .
