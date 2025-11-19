.PHONY: build
build: clean
	npx webpack

.PHONY: clean
clean:
	rm -rf dist

.PHONY: deps
deps:
	@npm install

.PHONY: lint
lint: deps
	@npx eslint js/**
