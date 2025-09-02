.PHONY: deps
deps:
	@pip3 install -r requirements.txt

.PHONY: test
test: deps
	@python3 -m pytest

.PHONY: lint
lint: deps
	@python3 -m pylint --recursive=true .
