.PHONY: run
run:
	@$(MAKE) -C server run

.PHONY: lint
lint:
	@$(MAKE) -C client lint
	@$(MAKE) -C server lint
	@$(MAKE) -C spec lint

PORT := 8080
export PORT

.PHONY: spec
spec:
	@$(MAKE) -C server run & \
		export SERVER_PID=$$!; \
		$(MAKE) -C spec test; \
		kill $${SERVER_PID}

.PHONY: pre-commit
pre-commit:
	pre-commit install
	pre-commit install --hook-type commit-msg
	pre-commit run --verbose --all-files --show-diff-on-failure
