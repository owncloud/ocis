
ifeq ($(MAKE_DEPTH),)
MAKE_DEPTH := 0
else
$(eval MAKE_DEPTH := $(shell echo "$$(( $(MAKE_DEPTH) + 1 ))" ) )
endif

export
