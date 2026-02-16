include ../.make/recursion.mk

############ tooling ############
ifneq (, $(shell command -v go 2> /dev/null)) # suppress `command not found warnings` for non go targets in CI
include ../.bingo/Variables.mk
endif

############ generate ############
include ../.make/generate.mk

.PHONY: ci-go-generate
ci-go-generate: $(MOCKERY) # CI runs ci-node-generate automatically before this target
	$(MOCKERY)
