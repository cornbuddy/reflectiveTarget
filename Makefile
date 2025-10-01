.PHONY: run
run:
	@$(MAKE) -C server run

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
	$(MAKE) -C server run & \
		$(MAKE) -C spec test && \
		pkill -f node

.PHONY: pre-commit
pre-commit:
	pre-commit install
	pre-commit install --hook-type commit-msg
	pre-commit run --verbose --all-files --show-diff-on-failure
