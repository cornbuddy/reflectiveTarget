LAYERS = app domain infra
MODULES = app domain infra test

.PHONY: all
all: lint test

.PHONY: lint
.ONESHELL:
lint:
	set -e
	for module in $(MODULES); do
		$(MAKE) -C $$module lint
	done

.PHONY: test
.ONESHELL:
test:
	set -e
	for layer in $(LAYERS); do
		ENVIRONMENT="development" $(MAKE) -C $$layer test
	done

.PHONY: deps
.ONESHELL:
deps:
	for module in $(MODULES); do
		$(MAKE) -C $$module deps
	done
