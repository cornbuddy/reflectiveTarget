include Makefile.venv

.PHONY: spec
spec: $(VENV)/pytest
	$(VENV)/pytest --verbose --html=report.html --self-contained-html $(ARGS)
	- xdg-open report.html

.PHONY: lint
lint: $(VENV)/pylint
	$(VENV)/pylint --recursive=true .
