BROWSER = firefox
REPORT = report.html

.PHONY: spec
spec: venv playwright
	$(VENV)/pytest \
		--dist=loadfile \
		--browser=$(BROWSER) \
		--html=$(REPORT) \
		$(if $(DEBUG),--headed --numprocesses=1,--numprocesses=auto) \
		$(ARGS)
	- xdg-open $(REPORT)

.PHONY: playwright
playwright: venv
	$(VENV)/playwright install $(BROWSER)

.PHONY: lint
lint: venv
	$(VENV)/pylint --recursive=true .

include Makefile.venv
