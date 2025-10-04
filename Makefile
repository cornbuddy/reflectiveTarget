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
		sleep 30 && \
		$(MAKE) -C spec test && \
		$(MAKE) -C server stop

.PHONY: pre-commit
pre-commit:
	pre-commit install
	pre-commit install --hook-type commit-msg
	pre-commit run --verbose --all-files --show-diff-on-failure
