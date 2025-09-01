.PHONY: pip
pip:
	@pip3 install -r requirements.txt

.PHONY: test
test: pip
	@python3 -m pytest

.PHONY: lint
lint: pip
	@python3 -m pylint --recursive=true .
