.PHONY: changelog
changelog: $(CALENS) ## generate changelog
	$(CALENS) -i ../changelog -t ../changelog/CHANGELOG.tmpl >| ../CHANGELOG.md

.PHONY: release
release: release-dirs release-linux release-windows release-darwin release-copy release-check

.PHONY: release-dirs
release-dirs:
	@mkdir -p $(DIST)/binaries $(DIST)/release

# docker specific packaging flags
DOCKER_LDFLAGS += -X "$(OCIS_REPO)/ocis-pkg/config/defaults.BaseDataPathType=path" -X "$(OCIS_REPO)/ocis-pkg/config/defaults.BaseDataPathValue=/var/lib/ocis"

release-linux-docker-amd64: release-dirs
	GOOS=linux \
	GOARCH=amd64 \
	go build \
		-tags 'netgo $(TAGS)' \
		-buildmode=pie \
		-trimpath \
		-ldflags '-extldflags "-static" $(LDFLAGS) $(DOCKER_LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-amd64-linux' \
		./cmd/$(NAME)

release-linux-docker-arm: release-dirs
	GOOS=linux \
	GOARCH=arm \
	go build \
		-tags 'netgo $(TAGS)' \
		-trimpath \
		-ldflags '-extldflags "-static" $(LDFLAGS) $(DOCKER_LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-arm-linux' \
		./cmd/$(NAME)

release-linux-docker-arm64: release-dirs
	GOOS=linux \
	GOARCH=arm64 \
	go build \
		-tags 'netgo $(TAGS)' \
		-buildmode=pie \
		-trimpath \
		-ldflags '-extldflags "-static" $(LDFLAGS) $(DOCKER_LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-arm64-linux' \
		./cmd/$(NAME)

.PHONY: release-linux
release-linux: release-dirs
	GOOS=linux \
	GOARCH=amd64 \
	go build \
		-tags 'netgo $(TAGS)' \
		-buildmode=pie \
		-trimpath \
		-ldflags '-extldflags "-static" $(LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-amd64-linux' \
		./cmd/$(NAME)

	GOOS=linux \
	GOARCH=386 \
	go build \
		-tags 'netgo $(TAGS)' \
		-buildmode=pie \
		-trimpath \
		-ldflags '-extldflags "-static" $(LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-386-linux' \
		./cmd/$(NAME)

	GOOS=linux \
	GOARCH=arm64 \
	go build \
		-tags 'netgo $(TAGS)' \
		-buildmode=pie \
		-trimpath \
		-ldflags '-extldflags "-static" $(LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-arm64-linux' \
		./cmd/$(NAME)

	GOOS=linux \
	GOARCH=arm \
	go build \
		-tags 'netgo $(TAGS)' \
		-trimpath \
		-ldflags '-extldflags "-static" $(LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-arm-linux' \
		./cmd/$(NAME)

.PHONY: release-windows
release-windows: release-dirs
	GOOS=windows \
	GOARCH=amd64 \
	go build \
		-tags 'netgo $(TAGS)' \
		-buildmode=pie \
		-trimpath \
		-ldflags '-extldflags "-static" $(LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-amd64-windows' \
		./cmd/$(NAME)

.PHONY: release-darwin
release-darwin: release-dirs
	GOOS=darwin \
	GOARCH=amd64 \
	go build \
		-tags 'netgo $(TAGS)' \
		-buildmode=pie \
		-trimpath \
		-ldflags '$(LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-amd64-darwin' \
		./cmd/$(NAME)

	GOOS=darwin \
	GOARCH=arm64 \
	go build \
		-tags 'netgo $(TAGS)' \
		-buildmode=pie \
		-trimpath \
		-ldflags '$(LDFLAGS)' \
		-o '$(DIST)/binaries/$(EXECUTABLE)-arm64-darwin' \
		./cmd/$(NAME)

.PHONY: release-copy
release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

.PHONY: release-check
release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/release/$(EXECUTABLE)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)

.PHONY: release-finish
release-finish: release-copy release-check
