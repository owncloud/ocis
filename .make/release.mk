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
release-linux-docker-amd64: $(GOX) release-dirs
	@$(GOX) -tags 'netgo $(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS) $(DOCKER_LDFLAGS)' -os 'linux' -arch 'amd64' -output '$(DIST)/binaries/$(EXECUTABLE)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

release-linux-docker-arm: $(GOX) release-dirs
	@$(GOX) -tags 'netgo $(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS) $(DOCKER_LDFLAGS)' -os 'linux' -arch 'arm' -output '$(DIST)/binaries/$(EXECUTABLE)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

release-linux-docker-arm64: $(GOX) release-dirs
	@$(GOX) -tags 'netgo $(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS) $(DOCKER_LDFLAGS)' -os 'linux' -arch 'arm64' -output '$(DIST)/binaries/$(EXECUTABLE)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-linux
release-linux: $(GOX) release-dirs
	@$(GOX) -tags 'netgo $(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS)' -os 'linux' -arch 'amd64 386 arm64 arm' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-windows
release-windows: $(GOX) release-dirs
	@$(GOX) -tags 'netgo $(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS)' -os 'windows' -arch 'amd64' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-darwin
release-darwin: $(GOX) release-dirs
	@$(GOX) -tags 'netgo $(TAGS)' -ldflags '$(LDFLAGS)' -os 'darwin' -arch 'amd64 arm64' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-copy
release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

.PHONY: release-check
release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/release/$(EXECUTABLE)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)

.PHONY: release-finish
release-finish: release-copy release-check
