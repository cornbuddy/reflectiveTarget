LAYERS = app domain infra

.PHONY: all
all: lint test

.PHONY: lint
lint:
	for layer in $(LAYERS) test; do \
		$(MAKE) -C $$layer lint; \
	done

.PHONY: test
test:
	for layer in $(LAYERS); do \
		ENVIRONMENT="development" $(MAKE) -C $$layer test; \
	done
