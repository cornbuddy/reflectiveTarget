.PHONY: run
run: stop
	docker compose -f develop/compose.yml up --build

.PHONY: stop
stop:
	- docker compose -f develop/compose.yml down

.PHONY: lint
lint:
	@$(MAKE) -C client lint
	@$(MAKE) -C server lint
	@$(MAKE) -C spec lint

.PHONY: test
test:
	@$(MAKE) -C server test

PORT := 8080
export PORT

.PHONY: spec
spec:
	$(MAKE) run & \
		sleep 15 \
		&& $(MAKE) -C spec spec \
		|| $(MAKE) stop

.PHONY: pre-commit
pre-commit:
	pre-commit install
	pre-commit install --hook-type commit-msg
	pre-commit run --verbose --all-files --show-diff-on-failure
