
.PHONY: docs-copy
docs-copy:
	mkdir -p hugo/content/; \
	cd hugo; \
	git init; \
	git remote rm origin; \
	git remote add origin https://github.com/owncloud/owncloud.github.io; \
	git fetch --depth=1; \
	git checkout origin/source -f; \
	rsync -ax --delete --exclude hugo/ --exclude Makefile ../. content/; \
