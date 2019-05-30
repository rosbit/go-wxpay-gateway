SHELL=/bin/bash

TAGS = gateway notify getsandbox

all:
	@for tag in $(TAGS); do \
		echo "building wxpay-$$tag ..."; \
		$(MAKE) -s -f make.inc s=static t=$$tag; \
	done

clean:
	@for tag in $(TAGS); do \
		rm -f go-wxpay-$$tag; \
	done
