.PHONY: all
all: test build

.PHONY: build
build: clean deps lint
	bun run build

.PHONY: clean
clean:
	- rm -rf dist node_modules

.PHONY: deps
deps:
	bun install
	bun update --latest
	bun audit fix --force

.PHONY: lint
lint: deps
	bunx eslint js/**

.PHONY: test
test: deps lint
	bun test $(ARGS)
