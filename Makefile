SHELL=/bin/bash

TAGS = gateway getsandbox
EXES = $(addprefix go-wxpay-,$(TAGS))

.PHONY: $(TAGS)
all: $(EXES)

go-wxpay-%: %
	@echo "building $@"
	$(MAKE) -s -f make.inc s=static t=$*

clean:
	@for tag in $(TAGS); do \
		rm -f go-wxpay-$$tag; \
	done
