include Makefile.venv

PIP ?= pip3

.PHONY: test
test: $(VENV)/pytest
	$(VENV)/pytest --verbose --html=report.html --self-contained-html \
		test/test_authz.py
	- xdg-open report.html

.PHONY: lint
lint: $(VENV)/pylint
	$(VENV)/pylint --recursive=true .
