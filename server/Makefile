LAYERS = app domain infra

.PHONY: lint
lint:
	for layer in $(LAYERS) test; do \
		$(MAKE) -C $$layer lint; \
	done

.PHONY: test
test:
	for layer in $(LAYERS); do \
		$(MAKE) -C $$layer test; \
	done
