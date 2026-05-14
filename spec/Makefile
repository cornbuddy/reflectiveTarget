.PHONY: spec
spec: venv
	$(VENV)/pytest --verbose --html=report.html --self-contained-html $(ARGS)
	- xdg-open report.html

.PHONY: lint
lint: venv
	$(VENV)/pylint --recursive=true .

include Makefile.venv
