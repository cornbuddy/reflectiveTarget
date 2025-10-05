include Makefile.venv

PIP ?= pip3

.PHONY: test
test: $(VENV)/pytest
	$(VENV)/pytest --html=report.html --self-contained-html
	- xdg-open report.html

.PHONY: lint
lint: $(VENV)/pylint
	$(VENV)/pylint --recursive=true .
