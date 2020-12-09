config-docs-generate:
	make -C ../accounts config-docs-generate
	make -C ../glauth config-docs-generate
	make -C ../konnectd config-docs-generate
	make -C ../ocis config-docs-generate
	make -C ../web config-docs-generate
	make -C ../ocis-pkg config-docs-generate
	make -C ../ocs config-docs-generate
	make -C ../proxy config-docs-generate
	make -C ../settings config-docs-generate
	make -C ../storage config-docs-generate
	make -C ../store config-docs-generate
	make -C ../thumbnails config-docs-generate
	make -C ../webdav config-docs-generate

.PHONY: docs-copy
docs-copy:
	mkdir -p hugo/content/
	cd hugo; \
	git init; \
	git config advice.detachedHead false; \
	git remote rm origin; \
	git remote add origin https://github.com/owncloud/owncloud.github.io; \
	git fetch --depth=1; \
	git checkout origin/source -f; \
	rsync -ax --delete --exclude hugo/ --exclude Makefile --exclude .gitignore --exclude README.md ../. content/; \

.PHONY: docs-serve
docs-serve: config-docs-generate docs-copy
	docker run --rm --network host -v $(shell pwd)/hugo:/src owncloudci/hugo:0 server

.PHONY: clean
clean:
	rm -rf hugo
