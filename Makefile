.PHONY: run
run: stop
	docker compose -f develop/compose.yml up --build

.PHONY: stop
stop:
	- docker compose -f develop/compose.yml down

.PHONY: lint
lint:
	@$(MAKE) -C client lint
	@$(MAKE) -C old-server lint
	@$(MAKE) -C spec lint

.PHONY: test
test:
	@$(MAKE) -C server test

PORT := 8080
export PORT

.PHONY: spec
spec:
	$(MAKE) -C old-server run & \
		sleep 30 && \
		$(MAKE) -C spec spec && \
		$(MAKE) -C old-server stop

.PHONY: pre-commit
pre-commit:
	pre-commit install
	pre-commit install --hook-type commit-msg
	pre-commit run --verbose --all-files --show-diff-on-failure
