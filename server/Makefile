.PHONY: run
run: stop deps
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

.PHONY: stop
stop:
	- pkill -f node
