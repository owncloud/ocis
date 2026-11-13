
ifeq ($(MAKE_DEPTH),)
MAKE_DEPTH := 0
else
$(eval MAKE_DEPTH := $(shell echo "$$(( $(MAKE_DEPTH) + 1 ))" ) )
endif

# Expand these variables one time, before "export". The bare "export" exports all variables, and
# GNU Make 4.4 also gives the exported variables to each $(shell ...). The variables of
# .bingo/Variables.mk run $(shell ...) on each use, so make ran these commands again and again.
# A value from the environment or from the command line has precedence.
GO         := $(or $(GO),$(shell which go))
GOPATH     := $(or $(GOPATH),$(shell $(GO) env GOPATH))
GOHOSTOS   := $(or $(GOHOSTOS),$(shell $(GO) env GOHOSTOS))
GOHOSTARCH := $(or $(GOHOSTARCH),$(shell $(GO) env GOHOSTARCH))
GOHOSTARM  := $(or $(GOHOSTARM),$(shell $(GO) env GOHOSTARM))

export
