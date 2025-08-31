.PHONY: pre-commit
pre-commit:
	pre-commit install
	pre-commit install --hook-type commit-msg
	pre-commit run --verbose --all-files --show-diff-on-failure

.PHONY: spec
spec:
	$(MAKE) -C server run & export SERVER_PID=$$!; \
		$(MAKE) -C spec test; \
		kill $${SERVER_PID}
