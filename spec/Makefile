include Makefile.venv

PIP ?= pip3

.PHONY: test
test: $(VENV)/pytest
	$(VENV)/pytest

.PHONY: lint
lint: $(VENV)/pylint
	$(VENV)/pylint --recursive=true .
