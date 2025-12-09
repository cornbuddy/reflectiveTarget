LAYERS = app domain infra
MODULES = app domain infra test

.PHONY: all
all: lint test

.PHONY: lint
lint:
	for module in $(MODULES); do \
		$(MAKE) -C $$module lint; \
	done

.PHONY: test
test:
	for layer in $(LAYERS); do \
		ENVIRONMENT="development" $(MAKE) -C $$layer test; \
	done

.PHONY: deps
deps:
	for module in $(MODULES); do \
		$(MAKE) -C $$module deps; \
	done
