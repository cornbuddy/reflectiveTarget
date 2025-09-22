.PHONY: run
run: deps
	@npm start

.PHONY: lint
lint: deps
	@npx eslint

.PHONY: test
test: deps lint
	@npm test

.PHONY: deps
deps:
	@npm install
