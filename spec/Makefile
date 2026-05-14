BROWSER = firefox
REPORT = report.html

define PYTEST
	$(VENV)/pytest \
		--browser=$(BROWSER) \
		--html=$(REPORT) \
		$(if $(DEBUG),--headed --numprocesses=1,--numprocesses=auto) \
		$(ARGS)
endef

.PHONY: spec
spec: venv playwright
	$(PYTEST)
	- xdg-open $(REPORT)

.PHONY: playwright
playwright: venv
	$(VENV)/playwright install $(BROWSER)

.PHONY: lint
lint: venv
	$(VENV)/pylint --recursive=true .

include Makefile.venv
