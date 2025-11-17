LAYERS = app model infra

.PHONY: tidy
tidy:
	for layer in $(LAYERS) test; do \
		$(MAKE) -C $$layer tidy; \
	done

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
