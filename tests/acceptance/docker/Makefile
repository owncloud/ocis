SHELL := bash

# define standard colors
BLACK        := $(shell tput -Txterm setaf 0)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
PURPLE       := $(shell tput -Txterm setaf 5)
BLUE         := $(shell tput -Txterm setaf 6)
WHITE        := $(shell tput -Txterm setaf 7)

RESET := $(shell tput -Txterm sgr0)

## default values only for sub-make calls
COMPOSE_FILE ?= src/redis.yml:src/ocis-base.yml:src/acceptance.yml

## user input
BEHAT_FEATURE ?=


ifdef OCIS_IMAGE_TAG
	BUILD_DEV_IMAGE := 0
else
	BUILD_DEV_IMAGE := 1
endif
OCIS_IMAGE_TAG ?= dev

# run tests with ociswrapper by default
WITH_WRAPPER ?= true
OCIS_WRAPPER := ../../ociswrapper/bin/ociswrapper

ifdef START_TIKA
	ifeq ($(START_TIKA),true)
		COMPOSE_FILE := $(COMPOSE_FILE):src/tika.yml
		SEARCH_EXTRACTOR_TYPE := tika
	else
		SEARCH_EXTRACTOR_TYPE := basic
	endif
else
	SEARCH_EXTRACTOR_TYPE := basic
endif

# static
DIVIDE_INTO_NUM_PARTS := 10
PARTS = 1 2 3 4 5 6 7 8 9 10
LOCAL_API_SUITES = $(shell ls ../features | grep ^api*)

COMPOSE_PROJECT_NAME := ocis-acceptance-tests

## make definition
.PHONY: help
help:
	@echo "Please use 'make <target>' where <target> is one of the following:"
	@echo
	@echo -e "${PURPLE}docs: https://owncloud.dev/ocis/development/testing/#testing-with-test-suite-in-docker${RESET}\n"
	@echo
	@echo -e "oCIS feature tests and test suites can be found here:"
	@echo -e "\thttps://github.com/owncloud/ocis/tree/master/tests/acceptance/features"
	@echo
	@echo -e "test suites that test ownCloud 10 API core compatibility are found here and they start with prefix coreApi-:"
	@echo -e "\thttps://github.com/owncloud/ocis/tree/master/tests/acceptance/features"
	@echo
	@echo -e "The oCIS to be tested will be build from your current working state."
	@echo -e "You also can select the oCIS Docker image for all tests by setting"
	@echo -e "\tmake ... ${YELLOW}OCIS_IMAGE_TAG=latest${RESET}"
	@echo -e "where ${YELLOW}latest${RESET} is an example for any valid Docker image tag from"
	@echo -e "https://hub.docker.com/r/owncloud/ocis."
	@echo
	@echo -e "${GREEN}Run full ownCloud test suites against oCIS with oCIS storage:${RESET}\n"
	@echo -e "\tmake Core-API-Tests-ocis-storage-${RED}X${RESET}\t\t${BLUE}run test suite number X, where ${RED}X = 1 .. 10${RESET}"
	@echo
	@echo -e "${GREEN}Run full ownCloud test suites against oCIS with s3ng storage:${RESET}\n"
	@echo -e "\tmake Core-API-Tests-s3ng-storage-${RED}X${RESET}\t\t${BLUE}run test suite number X, where ${RED}X = 1 .. 10${RESET}"
	@echo
	@echo -e "${GREEN}Run an oCIS feature test against oCIS with oCIS storage:${RESET}\n"
	@echo -e "\tmake test-ocis-feature-ocis-storage ${YELLOW}BEHAT_FEATURE='...'${RESET}\t${BLUE}run single feature test${RESET}"
	@echo
	@echo -e "${GREEN}Run an oCIS feature test against oCIS with s3ng storage:${RESET}\n"
	@echo -e "\tmake test-ocis-feature-s3ng-storage ${YELLOW}BEHAT_FEATURE='...'${RESET}\t${BLUE}run single feature test${RESET}"
	@echo
	@echo -e "${GREEN}Run an ownCloud feature test against oCIS with oCIS storage:${RESET}\n"
	@echo -e "\tmake test-core-feature-ocis-storage ${YELLOW}BEHAT_FEATURE='...'${RESET}\t${BLUE}run single feature test${RESET}"
	@echo
	@echo -e "\twhere ${YELLOW}BEHAT_FEATURE='...'${RESET} contains a relative path to the feature definition."
	@echo -e "\texample: ${RED}tests/acceptance/features/coreApiAuth/webDavAuth.feature${RESET}"
	@echo
	@echo -e "${GREEN}Run an ownCloud feature test against oCIS with s3ng storage:${RESET}\n"
	@echo -e "\tmake test-core-feature-s3ng-storage ${YELLOW}BEHAT_FEATURE='...'${RESET}\t${BLUE}run single feature test${RESET}"
	@echo
	@echo -e "\twhere ${YELLOW}BEHAT_FEATURE='...'${RESET} contains a relative path to the feature definition."
	@echo -e "\texample: ${RED}tests/acceptance/features/coreApiAuth/webDavAuth.feature${RESET}"
	@echo
	@echo -e "\twhere ${YELLOW}BEHAT_FEATURE='...'${RESET} contains a relative path to the feature definition."
	@echo -e "\texample: ${RED}tests/acceptance/features/coreApiAuth/webDavAuth.feature${RESET}"
	@echo
	@echo
	@echo -e "${GREEN}Show output of tests:${RESET}\n"
	@echo -e "\tmake show-test-logs\t\t${BLUE}show output of running or finished tests${RESET}"
	@echo
	@echo
	@echo -e "${GREEN}Clean up after testing:${RESET}\n"
	@echo -e "\tmake clean\t${BLUE}clean up all${RESET}"
	@echo -e "\tmake clean-docker-container\t\t${BLUE}stops and removes used docker containers${RESET}"
	@echo -e "\tmake clean-docker-volumes\t\t${BLUE}removes used docker volumes (used for caching)${RESET}"
	@echo
.PHONY: test-ocis-feature-ocis-storage
test-ocis-feature-ocis-storage:
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=ocis \
	BEHAT_FEATURE=$(BEHAT_FEATURE) \
	$(MAKE) --no-print-directory testSuite

.PHONY: test-ocis-feature-s3ng-storage
test-ocis-feature-s3ng-storage:
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=s3ng \
	BEHAT_FEATURE=$(BEHAT_FEATURE) \
	START_CEPH=1 \
	$(MAKE) --no-print-directory testSuite

.PHONY: test-core-feature-ocis-storage
test-core-feature-ocis-storage: ## test a core feature with oCIS storage, usage: make ... BEHAT_FEATURE='tests/acceptance/features/coreApiAuth/webDavAuth.feature'
	@TEST_SOURCE=oc10 \
	STORAGE_DRIVER=ocis \
	BEHAT_FEATURE=$(BEHAT_FEATURE) \
	$(MAKE) --no-print-directory testSuite

.PHONY: test-core-feature-s3ng-storage
test-core-feature-s3ng-storage: ## test a core feature with s3ng storage, usage: make ... BEHAT_FEATURE='tests/acceptance/features/coreApiAuth/webDavAuth.feature'
	@TEST_SOURCE=oc10 \
	STORAGE_DRIVER=s3ng \
	BEHAT_FEATURE=$(BEHAT_FEATURE) \
	START_CEPH=1 \
	$(MAKE) --no-print-directory testSuite

localSuiteOcis = $(addprefix localApiTests-, $(addsuffix -ocis,${LOCAL_API_SUITES}))
.PHONY: $(localSuiteOcis)
$(localSuiteOcis): ## run local api test suite with ocis storage
	@$(eval BEHAT_SUITE=$(shell echo "$@" | cut -d'-' -f2))
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=ocis \
	BEHAT_SUITE=$(BEHAT_SUITE) \
	$(MAKE) --no-print-directory testSuite

localSuites3ng = $(addprefix localApiTests-, $(addsuffix -s3ng,${LOCAL_API_SUITES}))
.PHONY: $(localSuites3ng)
$(localSuites3ng): ## run local api test suite with s3 storage
	@$(eval BEHAT_SUITE=$(shell echo "$@" | cut -d'-' -f2))
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=s3ng \
	BEHAT_SUITE=$(BEHAT_SUITE) \
	$(MAKE) --no-print-directory testSuite

targetsOCIS = $(addprefix Core-API-Tests-ocis-storage-,$(PARTS))
.PHONY: $(targetsOCIS)
$(targetsOCIS):
	@$(eval RUN_PART=$(shell echo "$@" | tr -dc '0-9'))
	@TEST_SOURCE=oc10 \
	STORAGE_DRIVER=ocis \
	RUN_PART=$(RUN_PART) \
	$(MAKE) --no-print-directory testSuite

targetsS3ng = $(addprefix Core-API-Tests-s3ng-storage-,$(PARTS))
.PHONY: $(targetsS3ng)
$(targets):
	@$(eval RUN_PART=$(shell echo "$@" | tr -dc '0-9'))
	@TEST_SOURCE=oc10 \
	STORAGE_DRIVER=s3ng \
	RUN_PART=$(RUN_PART) \
	$(MAKE) --no-print-directory testSuite

.PHONY: testSuite
testSuite: $(OCIS_WRAPPER) build-dev-image clean-docker-container
	@if [ -n "${START_CEPH}" ]; then \
		COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
		COMPOSE_FILE=src/ceph.yml \
		docker-compose run start_ceph; \
	fi; \
	if [ "${START_EMAIL}" == "true" ]; then \
		COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
		COMPOSE_FILE=src/email.yml \
		docker-compose run start_email; \
	fi; \
	if [ "${START_ANTIVIRUS}" == "true" ]; then \
    		COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
    		COMPOSE_FILE=src/antivirus.yml \
    		docker-compose run start_antivirus; \
    fi; \
    if [ "${START_TIKA}" == "true" ]; then \
			COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
			COMPOSE_FILE=src/tika.yml \
			docker-compose run tika-service; \
    fi; \
	COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
	COMPOSE_FILE=$(COMPOSE_FILE) \
	STORAGE_DRIVER=$(STORAGE_DRIVER) \
	TEST_SOURCE=$(TEST_SOURCE) \
	WITH_WRAPPER=$(WITH_WRAPPER) \
	OCIS_ASYNC_UPLOADS=$(OCIS_ASYNC_UPLOADS) \
	OCIS_ADD_RUN_SERVICES=$(OCIS_ADD_RUN_SERVICES) \
	POSTPROCESSING_STEPS=$(POSTPROCESSING_STEPS) \
	SEARCH_EXTRACTOR_TYPE=$(SEARCH_EXTRACTOR_TYPE) \
	OCIS_IMAGE_TAG=$(OCIS_IMAGE_TAG) \
	BEHAT_SUITE=$(BEHAT_SUITE) \
	BEHAT_FEATURE=$(BEHAT_FEATURE) \
	DIVIDE_INTO_NUM_PARTS=$(DIVIDE_INTO_NUM_PARTS) \
	RUN_PART=$(RUN_PART) \
	docker-compose up -d --build --force-recreate

.PHONY: show-test-logs
show-test-logs: ## show logs of test
	@COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
	COMPOSE_FILE=$(COMPOSE_FILE) \
	docker-compose logs --no-log-prefix -f acceptance-tests | less

$(OCIS_WRAPPER):
	@if [ "$(WITH_WRAPPER)" == "true" ]; then \
		$(MAKE) --no-print-directory -C ../../ociswrapper build \
	; fi;

.PHONY: build-dev-image
build-dev-image:
	@if [ $(BUILD_DEV_IMAGE) -eq 1 ]; then \
		$(MAKE) --no-print-directory -C ../../../ocis dev-docker \
	; fi;

.PHONY: clean-dev-docker-image
clean-dev-docker-image: ## clean docker image built during acceptance tests
	@docker image rm owncloud/ocis:dev || true

.PHONY: clean-docker-container
clean-docker-container: ## clean docker containers created during acceptance tests
	@COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
	COMPOSE_FILE=$(COMPOSE_FILE) \
	BEHAT_SUITE="" \
	DIVIDE_INTO_NUM_PARTS="" \
	OCIS_IMAGE_TAG="" \
	RUN_PART="" \
	STORAGE_DRIVER="" \
	TEST_SOURCE="" \
	docker-compose down --remove-orphans

.PHONY: clean-docker-volumes
clean-docker-volumes: ## clean docker volumes created during acceptance tests
	@COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
	COMPOSE_FILE=$(COMPOSE_FILE) \
	BEHAT_SUITE="" \
	DIVIDE_INTO_NUM_PARTS="" \
	OCIS_IMAGE_TAG="" \
	RUN_PART="" \
	STORAGE_DRIVER="" \
	TEST_SOURCE="" \
	docker-compose down --remove-orphans -v

.PHONY: clean-files
clean-files:
	@$(MAKE) --no-print-directory -C ../../../. clean-tests

.PHONY: clean
clean: clean-docker-container clean-docker-volumes clean-dev-docker-image clean-files ## clean all
