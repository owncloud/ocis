.PHONY: changelog
changelog: $(CALENS)
	$(CALENS) >| CHANGELOG.md