.PHONY: all
all: build test

.PHONY: build
build: deps lint clean
	@npx webpack

.PHONY: clean
clean:
	- rm -rf dist

.PHONY: deps
deps:
	@npm install

.PHONY: lint
lint: deps
	@npx eslint js/**

.PHONY: test
test: deps lint
	@npx jest --detectOpenHandles $(ARGS)
