SHELL := bash
NAME := ocis
IMPORT := github.com/owncloud/$(NAME)
BIN := bin
DIST := dist
HUGO := hugo
CONFIG := config/identifier-registration.yaml

ifeq ($(OS), Windows_NT)
	EXECUTABLE := $(NAME).exe
	UNAME := Windows
else
	EXECUTABLE := $(NAME)
	UNAME := $(shell uname -s)
endif

ifeq ($(UNAME), Darwin)
	GOBUILD ?= go build -i
else
	GOBUILD ?= go build
endif

PACKAGES ?= $(shell go list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f)
GENERATE ?= $(PACKAGES)

TAGS ?=

ifndef OUTPUT
	ifneq ($(DRONE_TAG),)
		OUTPUT ?= $(subst v,,$(DRONE_TAG))
	else
		OUTPUT ?= testing
	endif
endif

ifndef VERSION
	ifneq ($(DRONE_TAG),)
		VERSION ?= $(subst v,,$(DRONE_TAG))
	else
		VERSION ?= $(shell git rev-parse --short HEAD)
	endif
endif

ifndef DATE
	DATE := $(shell date -u '+%Y%m%d')
endif

LDFLAGS += -s -w -X "$(IMPORT)/pkg/version.String=$(VERSION)" -X "$(IMPORT)/pkg/version.Date=$(DATE)"
DEBUG_LDFLAGS += -X "$(IMPORT)/pkg/version.String=$(VERSION)" -X "$(IMPORT)/pkg/version.Date=$(DATE)"
GCFLAGS += all=-N -l

.PHONY: all
all: build

.PHONY: sync
sync:
	go mod download

.PHONY: clean
clean: clean-config
	go clean -i ./...
	rm -rf $(BIN) $(DIST) $(HUGO)

.PHONY: clean-config
clean-config:
	rm -rf $(CONFIG)

.PHONY: fmt
fmt:
	gofmt -s -w $(SOURCES)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: lint
lint:
	for PKG in $(PACKAGES); do go run golang.org/x/lint/golint -set_exit_status $$PKG || exit 1; done;

.PHONY: generate
generate:
	go generate $(GENERATE)

.PHONY: changelog
changelog:
	go run github.com/restic/calens >| CHANGELOG.md

.PHONY: test
test:
	go run github.com/haya14busa/goverage -v -coverprofile coverage.out $(PACKAGES)

.PHONY: install
install: $(SOURCES)
	go install -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' ./cmd/$(NAME)

.PHONY: build
build: $(BIN)/$(EXECUTABLE) $(BIN)/$(EXECUTABLE)-debug

$(BIN)/$(EXECUTABLE): $(SOURCES)
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $@ ./cmd/$(NAME)

$(BIN)/$(EXECUTABLE)-debug: $(SOURCES)
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(DEBUG_LDFLAGS)' -gcflags '$(GCFLAGS)' -o $@ ./cmd/$(NAME)

$(BIN)/$(EXECUTABLE)-linux: $(SOURCES)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -gcflags '$(GCFLAGS)' -o $@ ./cmd/$(NAME)

.PHONY: staticcheck
staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck -tags '$(TAGS)' $(PACKAGES)

.PHONY: release
release: release-dirs release-linux release-windows release-darwin release-copy release-check

.PHONY: release-dirs
release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

.PHONY: release-linux
release-linux: release-dirs
	go run github.com/mitchellh/gox -tags 'netgo $(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS)' -os 'linux' -arch 'amd64 386 arm64 arm' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-windows
release-windows: release-dirs
	go run github.com/mitchellh/gox -tags 'netgo $(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS)' -os 'windows' -arch 'amd64' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-darwin
release-darwin: release-dirs
	go run github.com/mitchellh/gox -tags 'netgo $(TAGS)' -ldflags '$(LDFLAGS)' -os 'darwin' -arch 'amd64' -output '$(DIST)/binaries/$(EXECUTABLE)-$(OUTPUT)-{{.OS}}-{{.Arch}}' ./cmd/$(NAME)

.PHONY: release-copy
release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

.PHONY: release-check
release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/release/$(EXECUTABLE)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)

.PHONY: release-finish
release-finish: release-copy release-check

.PHONY: docs-copy
docs-copy:
	mkdir -p $(HUGO); \
	mkdir -p $(HUGO)/content/; \
	cd $(HUGO); \
	git init; \
	git remote rm origin; \
	git remote add origin https://github.com/owncloud/owncloud.github.io; \
	git fetch --depth=1; \
	git checkout origin/source -f; \
	rsync --delete -ax --exclude 'static' ../docs/ content/$(NAME); \
	rsync --delete -ax ../docs/static/ static/$(NAME); \

.PHONY: config-docs-generate
config-docs-generate:
	go run github.com/owncloud/flaex >| docs/configuration.md

.PHONY: docs-build
docs-build:
	cd $(HUGO); hugo

.PHONY: docs
docs: config-docs-generate docs-copy docs-build

.PHONY: watch
watch:
	go run github.com/cespare/reflex -c reflex.conf

# -------------------------------------------------------------------------------
# EOS related destinations
# -------------------------------------------------------------------------------

EOS_LDAP_HOST ?= host.docker.internal:9125

eos-docker:
	git clone https://gitlab.cern.ch/eos/eos-docker.git

eos-docker/scripts/start_services_ocis.sh: eos-docker
	# TODO find a way to properly inject the following env vars into the container:
	# EOS_UTF8=1 enables utf8 filenames
	# EOS_NS_ACCOUNTING=1 enables dir size propagation
	# EOS_SYNCTIME_ACCOUNTING=1 enables mtime propagation
	#  - needs the sys.mtime.propagation=1 on a home dir, handled by the reva eos storage driver
	#  - sys.allow.oc.sync=1 is not needed, it is an option for the eos built in webdav endpoint
	# 1. -e: for now, we patch the start_services.sh and use that
	# 2. -e: we need to expose the storageprovider ports whan running the docker containen
	# TODO use port from address to open different ports, this currently only works for one client container
	sed -e "s/--name eos-mgm1 --net/--name eos-mgm1 --env EOS_UTF8=1 --env EOS_NS_ACCOUNTING=1 --env EOS_SYNCTIME_ACCOUNTING=1 --net/" -e 's/--name $${CLIENTHOSTNAME} --net=eoscluster.cern.ch/--name $${CLIENTHOSTNAME} -p 9154:9154 -p 9155:9155 -p 9156:9156 -p 9157:9157 -p 9158:9158 -p 9159:9159 -p 9160:9160 -p 9161:9161 --net=eoscluster.cern.ch/' ./eos-docker/scripts/start_services.sh > ./eos-docker/scripts/start_services_ocis.sh
	chmod +x ./eos-docker/scripts/start_services_ocis.sh

.PHONY: eos-deploy
eos-deploy: eos-docker/scripts/start_services_ocis.sh
	# TODO keep eos up to date: see https://gitlab.cern.ch/dss/eos/tags
	./eos-docker/scripts/start_services_ocis.sh -i gitlab-registry.cern.ch/dss/eos:4.7.12 -q
	# Install ldap packages
	docker exec -i eos-mgm1 yum install -y nss-pam-ldapd nscd authconfig
	docker exec -i eos-cli1 yum install -y nss-pam-ldapd nscd authconfig

.PHONY: eos-setup
eos-setup: eos-docker/scripts/start_services_ocis.sh
	#Allow resolving uids against ldap
	# 9125 is the ldap port, 9126 would be tls ... but self signed cert
	# TODO check out the error message (ignoring for now ... still works): read LDAP host from env var, if not set fall back to docker host, in docker compose should be the ocis-glauth container because it contains guest accounts a well
ifeq ($(UNAME), Linux)
	#on linux add host.docker.internal to hosts: https://stackoverflow.com/questions/714100/os-detecting-makefile
	docker exec -it eos-mgm1 /bin/sh -c $$'echo -e "`/sbin/ip route | awk \'/default/ { print $$3 }\'`\thost.docker.internal" | sudo tee -a /etc/hosts > /dev/null'
	docker exec -it eos-cli1 /bin/sh -c $$'echo -e "`/sbin/ip route | awk \'/default/ { print $$3 }\'`\thost.docker.internal" | sudo tee -a /etc/hosts > /dev/null'
endif
	docker exec -i eos-mgm1 authconfig --enableldap --enableldapauth --ldapserver=$(EOS_LDAP_HOST) --ldapbasedn="dc=example,dc=org" --update; \
	docker exec -i eos-cli1 authconfig --enableldap --enableldapauth --ldapserver=$(EOS_LDAP_HOST) --ldapbasedn="dc=example,dc=org" --update;

	# setup users on mgm
	#TODO Failed to get D-Bus connection: Operation not permitted\ngetsebool:  SELinux is disabled
	docker exec -i eos-mgm1 sed -i "s/#binddn cn=.*/binddn cn=reva,ou=sysusers,dc=example,dc=org/" /etc/nslcd.conf
	docker exec -i eos-mgm1 sed -i "s/#bindpw .*/bindpw reva/" /etc/nslcd.conf
	# print the actual authconfig
	docker exec -i eos-mgm1 authconfig --test
	# start nslcd. you need to restart it if you change the ldap config
	docker exec -i eos-mgm1 nslcd
	# use unix accounts
	docker exec -i eos-mgm1 eos vid set map -unix "<pwd>" vuid:0 vgid:0
	# allow cli to create homes
	docker exec -i eos-mgm1 eos vid add gateway eos-cli1
	# krb not needed
	docker exec -i eos-mgm1 eos vid disable krb5

	# setup users on cli, same as for mgm
	docker exec -i eos-cli1 sed -i "s/#binddn cn=.*/binddn cn=reva,ou=sysusers,dc=example,dc=org/" /etc/nslcd.conf
	docker exec -i eos-cli1 sed -i "s/#bindpw .*/bindpw reva/" /etc/nslcd.conf
	docker exec -i eos-cli1 nslcd

	# create necessary lib link for ocis
	docker exec -i eos-cli1 ln -s /lib64/ld-linux-x86-64.so.2 /lib

.PHONY: eos-test
eos-test:
	# check we know the demo users
	docker exec -i eos-mgm1 id einstein
	docker exec -i eos-mgm1 id marie
	docker exec -i eos-mgm1 id feynman

.PHONY: eos-copy-ocis
eos-copy-ocis: build $(BIN)/$(EXECUTABLE)-linux
	# copy the linux binary to the eos-cli1 container
	docker cp ./bin/ocis-linux eos-cli1:/usr/local/bin/ocis

.PHONY: eos-ocis-storage-home
eos-ocis-storage-home:
	# configure the home storage to use the eos driver and return the mount id of the eos driver in responses
	# mount a set of eoshome storage drivers for requests to /webdav
	docker exec -i \
	--env OCIS_LOG_LEVEL=debug \
	--env REVA_STORAGE_HOME_DRIVER=eoshome \
	--env REVA_STORAGE_HOME_MOUNT_ID=1284d238-aa92-42ce-bdc4-0b0000009154 \
	eos-cli1 ocis reva-storage-home &
	docker exec -i \
	--env OCIS_LOG_LEVEL=debug \
	--env REVA_STORAGE_HOME_DATA_DRIVER=eoshome \
	--env REVA_GATEWAY_URL=host.docker.internal:9142 \
	eos-cli1 ocis reva-storage-home-data &
	# mount a second set of eoshome storage drivers for requests to /dav/files
	docker exec -i \
	--env OCIS_LOG_LEVEL=debug \
	--env REVA_STORAGE_EOS_DRIVER=eoshome \
	--env REVA_STORAGE_EOS_NAMESPACE="/eos/dockertest/reva/users" \
	--env REVA_STORAGE_EOS_LAYOUT="{{substr 0 1 .Username}}" \
	eos-cli1 ocis reva-storage-eos &
	docker exec -i \
	--env OCIS_LOG_LEVEL=debug \
	--env REVA_STORAGE_EOS_DATA_DRIVER=eoshome \
	--env REVA_STORAGE_EOS_NAMESPACE="/eos/dockertest/reva/users" \
	--env REVA_STORAGE_EOS_LAYOUT="{{substr 0 1 .Username}}" \
	--env REVA_GATEWAY_URL=host.docker.internal:9142 \
	eos-cli1 ocis reva-storage-eos-data &

.PHONY: eos-ocis
eos-ocis:
	export OCIS_LOG_LEVEL=debug; \
	export DAV_FILES_NAMESPACE="/eos/"; \
	bin/ocis micro & \
	bin/ocis glauth & \
	bin/ocis graph-explorer & \
	bin/ocis graph & \
	bin/ocis konnectd & \
	bin/ocis phoenix & \
	bin/ocis thumbnails & \
	bin/ocis webdav & \
	bin/ocis reva-auth-basic & \
	bin/ocis reva-auth-bearer & \
	bin/ocis reva-frontend & \
	bin/ocis reva-storage-public-link & \
	bin/ocis reva-gateway & \
	bin/ocis reva-sharing & \
	bin/ocis reva-users & \
	bin/ocis proxy &

.PHONY: eos-start
eos-start: eos-deploy eos-setup eos-copy-ocis eos-ocis-storage-home eos-ocis

.PHONY: eos-clean
eos-clean:
	rm eos-docker/scripts/start_services_ocis.sh

.PHONY: eos-stop
eos-stop: eos-docker
	./eos-docker/scripts/shutdown_services.sh

.PHONY: eos-install-go
eos-install-go:
	docker exec -i eos-cli1 curl https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz -O
	docker exec -i eos-cli1 tar -C /usr/local -xzf go1.14.4.linux-amd64.tar.gz
	# export PATH=$PATH:/usr/local/go/bin
