COMPOSE = docker compose -f develop/compose.yml

.PHONY: all
all: build run

.PHONY: run
run: stop
	$(COMPOSE) up

.PHONY: stop
stop:
	- $(COMPOSE) down

.PHONY: clean
clean: stop
	$(COMPOSE) down --volumes

.PHONY: build
build:
	$(COMPOSE) pull --include-deps
	$(COMPOSE) build --pull $(BUILD_ARGS)

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
	@$(MAKE) -C client test
	@$(MAKE) -C server test

.PHONY: spec
spec: clean build
	$(MAKE) run &
	sleep 10
	$(MAKE) -C spec spec; ret=$$?; \
		$(MAKE) stop; \
		exit $$ret

.PHONY: pre-commit
pre-commit:
	pre-commit install
	pre-commit install --hook-type commit-msg
	pre-commit run --verbose --all-files --show-diff-on-failure

.PHONY: venv
venv:
	$(eval OS_SHELL := $(shell basename $$SHELL))
	$(MAKE) -C spec $(OS_SHELL)
