.PHONY: run
run: stop
	docker compose -f ../develop/compose.yml up --build

.PHONY: lint
lint: deps
	@npx eslint

.PHONY: test
test: deps lint
	@npx jest --detectOpenHandles $(ARGS)

.PHONY: deps
deps:
	@npm install

.PHONY: stop
stop:
	- docker compose -f ../develop/compose.yml down
