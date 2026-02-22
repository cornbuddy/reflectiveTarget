COMPOSE = develop/compose.yml

.PHONY: all
all: build run

.PHONY: run
run: stop
	docker compose -f $(COMPOSE) up

.PHONY: stop
stop:
	- docker compose -f $(COMPOSE) down

.PHONY: build
build:
	docker compose -f $(COMPOSE) build --pull

.PHONY: lint
lint:
	@$(MAKE) -C client lint
	@$(MAKE) -C server lint
	@$(MAKE) -C spec lint

.PHONY: deps
deps:
	@$(MAKE) -C client deps
	@$(MAKE) -C server deps

.PHONY: test
test:
	@$(MAKE) -C server test

.PHONY: spec
spec: build
	$(MAKE) run &
	sleep 10
	- $(MAKE) -C spec spec
	$(MAKE) stop

.PHONY: pre-commit
pre-commit:
	pre-commit install
	pre-commit install --hook-type commit-msg
	pre-commit run --verbose --all-files --show-diff-on-failure
