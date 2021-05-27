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
	COMPOSE_FILE := $(COMPOSE_FILE):src/ocis-image.yml
else
	COMPOSE_FILE := $(COMPOSE_FILE):src/ocis-build.yml
endif
OCIS_IMAGE_TAG ?= latest

# static
DIVIDE_INTO_NUM_PARTS := 10
PARTS = 1 2 3 4 5 6 7 8 9 10

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
	@echo -e "ownCloud feature tests and test suites can be found here:"
	@echo -e "\thttps://github.com/owncloud/core/tree/master/tests/acceptance/features"
	@echo
	@echo -e "The oCIS to be tested will be build from your current working state."
	@echo -e "You also can select the oCIS Docker image for all tests by setting"
	@echo -e "\tmake ... ${YELLOW}OCIS_IMAGE_TAG=latest${RESET}"
	@echo -e "where ${YELLOW}latest${RESET} is an example for any valid Docker image tag from"
	@echo -e "https://hub.docker.com/r/owncloud/ocis."
	@echo
	@echo -e "${GREEN}Run full oCIS test suites against oCIS with oCIS storage:${RESET}\n"
	@echo -e "\tmake localApiTests-apiAccountsHashDifficulty-ocis\t\t${BLUE}run apiAccountsHashDifficulty test suite${RESET}"
	@echo -e "\tmake localApiTests-apiBugDemonstration-ocis\t\t${BLUE}run apiBugDemonstration test suite${RESET}"
	@echo
	@echo -e "${GREEN}Run full oCIS test suites against oCIS with s3ng storage:${RESET}\n"
	@echo -e "\tmake localApiTests-apiBasic-s3ng\t\t${BLUE}run apiBasic test suite${RESET}"
	@echo -e "\tmake localApiTests-apiOcisSpecific-s3ng\t\t${BLUE}run apiOcisSPecific test suite${RESET}"
	@echo
	@echo -e "${GREEN}Run full oCIS test suites against oCIS with ownCloud storage:${RESET}\n"
	@echo -e "\tmake localApiTests-apiAccountsHashDifficulty-owncloud\t\t${BLUE}run apiAccountsHashDifficulty test suite${RESET}"
	@echo -e "\tmake localApiTests-apiBugDemonstration-owncloud\t${BLUE}run apiBugDemonstration test suite${RESET}"
	@echo
	@echo -e "${GREEN}Run full ownCloud test suites against oCIS with oCIS storage:${RESET}\n"
	@echo -e "\tmake Core-API-Tests-ocis-storage-${RED}X${RESET}\t\t${BLUE}run test suite number X, where ${RED}X = 1 .. 10${RESET}"
	@echo
	@echo -e "${GREEN}Run full ownCloud test suites against oCIS with s3ng storage:${RESET}\n"
	@echo -e "\tmake Core-API-Tests-s3ng-storage-${RED}X${RESET}\t\t${BLUE}run test suite number X, where ${RED}X = 1 .. 10${RESET}"
	@echo
	@echo -e "${GREEN}Run full ownCloud test suites against oCIS with ownCloud storage:${RESET}\n"
	@echo -e "\tmake Core-API-Tests-owncloud-storage-${RED}X${RESET}\t\t${BLUE}run test suite number X, where ${RED}X = 1 .. 10${RESET}"
	@echo
	@echo -e "${GREEN}Run an oCIS feature test against oCIS with oCIS storage:${RESET}\n"
	@echo -e "\tmake test-ocis-feature-ocis-storage ${YELLOW}BEHAT_FEATURE='...'${RESET}\t${BLUE}run single feature test${RESET}"
	@echo
	@echo -e "\twhere ${YELLOW}BEHAT_FEATURE='...'${RESET} contains a relative path to the feature definition."
	@echo -e "\texample: ${RED}tests/acceptance/features/apiBugDemonstration/apiAuthOcs-ocsDELETEAuth.feature${RESET}"
	@echo
	@echo -e "${GREEN}Run an oCIS feature test against oCIS with s3ng storage:${RESET}\n"
	@echo -e "\tmake test-ocis-feature-s3ng-storage ${YELLOW}BEHAT_FEATURE='...'${RESET}\t${BLUE}run single feature test${RESET}"
	@echo
	@echo -e "\twhere ${YELLOW}BEHAT_FEATURE='...'${RESET} contains a relative path to the feature definition."
	@echo -e "\texample: ${RED}tests/acceptance/features/apiOcisSpecific/apiAuthOcs-ocsDELETEAuth.feature${RESET}"
	@echo
	@echo -e "${GREEN}Run an oCIS feature test against oCIS with owncloud storage:${RESET}\n"
	@echo -e "\tmake test-ocis-feature-owncloud-storage ${YELLOW}BEHAT_FEATURE='...'${RESET}\t${BLUE}run single feature test${RESET}"
	@echo
	@echo -e "\twhere ${YELLOW}BEHAT_FEATURE='...'${RESET} contains a relative path to the feature definition."
	@echo -e "\texample: ${RED}tests/acceptance/features/apiBugDemonstration/apiAuthOcs-ocsDELETEAuth.feature${RESET}"
	@echo
	@echo -e "${GREEN}Run an ownCloud feature test against oCIS with oCIS storage:${RESET}\n"
	@echo -e "\tmake test-oc10-feature-ocis-storage ${YELLOW}BEHAT_FEATURE='...'${RESET}\t${BLUE}run single feature test${RESET}"
	@echo
	@echo -e "\twhere ${YELLOW}BEHAT_FEATURE='...'${RESET} contains a relative path to the feature definition."
	@echo -e "\texample: ${RED}tests/acceptance/features/apiAuth/cors.feature${RESET}"
	@echo
	@echo -e "${GREEN}Run an ownCloud feature test against oCIS with s3ng storage:${RESET}\n"
	@echo -e "\tmake test-oc10-feature-s3ng-storage ${YELLOW}BEHAT_FEATURE='...'${RESET}\t${BLUE}run single feature test${RESET}"
	@echo
	@echo -e "\twhere ${YELLOW}BEHAT_FEATURE='...'${RESET} contains a relative path to the feature definition."
	@echo -e "\texample: ${RED}tests/acceptance/features/apiAuth/cors.feature${RESET}"
	@echo
	@echo -e "${GREEN}Run an ownCloud feature test against oCIS with owncloud storage:${RESET}\n"
	@echo -e "\tmake test-oc10-feature-owncloud-storage ${YELLOW}BEHAT_FEATURE='...'${RESET}\t${BLUE}run single feature test${RESET}"
	@echo
	@echo -e "\twhere ${YELLOW}BEHAT_FEATURE='...'${RESET} contains a relative path to the feature definition."
	@echo -e "\texample: ${RED}tests/acceptance/features/apiAuth/cors.feature${RESET}"
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
test-ocis-feature-ocis-storage: ## test a ocis feature with oCIS storage, useage: make ... BEHAT_FEATURE='tests/acceptance/features/apiBugDemonstration/apiAuthOcs-ocsDELETEAuth.feature:7'
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=ocis \
	BEHAT_FEATURE=$(BEHAT_FEATURE) \
	$(MAKE) --no-print-directory testSuite

.PHONY: test-ocis-feature-s3ng-storage
test-ocis-feature-s3ng-storage: ## test a ocis feature with s3ng storage, useage: make ... BEHAT_FEATURE='tests/acceptance/features/apiOcisSpecific/apiAuthOcs-ocsDELETEAuth.feature:7'
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=s3ng \
	BEHAT_FEATURE=$(BEHAT_FEATURE) \
	START_CEPH=1 \
	$(MAKE) --no-print-directory testSuite

.PHONY: test-ocis-feature-owncloud-storage
test-ocis-feature-owncloud-storage: ## test a ocis feature with oc10 storage, useage: make ... BEHAT_FEATURE='tests/acceptance/features/apiBugDemonstration/apiAuthOcs-ocsDELETEAuth.feature:7'
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=owncloud \
	BEHAT_FEATURE=$(BEHAT_FEATURE) \
	$(MAKE) --no-print-directory testSuite

.PHONY: test-oc10-feature-ocis-storage
test-oc10-feature-ocis-storage: ## test a oC10 feature with oCIS storage, useage: make ... BEHAT_FEATURE='tests/acceptance/features/apiAuth/cors.feature'
	@TEST_SOURCE=oc10 \
	STORAGE_DRIVER=ocis \
	BEHAT_FEATURE=$(BEHAT_FEATURE) \
	$(MAKE) --no-print-directory testSuite

.PHONY: test-oc10-feature-s3ng-storage
test-oc10-feature-s3ng-storage: ## test a oC10 feature with s3ng storage, useage: make ... BEHAT_FEATURE='tests/acceptance/features/apiAuth/cors.feature'
	@TEST_SOURCE=oc10 \
	STORAGE_DRIVER=s3ng \
	BEHAT_FEATURE=$(BEHAT_FEATURE) \
	START_CEPH=1 \
	$(MAKE) --no-print-directory testSuite

.PHONY: test-oc10-feature-owncloud-storage
test-oc10-feature-owncloud-storage: ## test a oC10 feature with oc10 storage, useage: make ... BEHAT_FEATURE='tests/acceptance/features/apiAuth/cors.feature'
	@TEST_SOURCE=oc10 \
	STORAGE_DRIVER=owncloud \
	BEHAT_FEATURE=$(BEHAT_FEATURE) \
	$(MAKE) --no-print-directory testSuite

.PHONY: localApiTests-apiBugDemonstration-owncloud
localApiTests-apiBugDemonstration-owncloud: ## run apiBugDemonstration test suite with owncloud storage
	@TEST_SOURCE=oc10 \
	STORAGE_DRIVER=owncloud \
	BEHAT_SUITE=apiBugDemonstration \
	$(MAKE) --no-print-directory testSuite

.PHONY: localApiTests-apiAccountsHashDifficulty-owncloud
localApiTests-apiAccountsHashDifficulty-owncloud: ## run apiAccountsHashDifficulty test suite with owncloud storage
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=owncloud \
	BEHAT_SUITE=apiAccountsHashDifficulty \
	$(MAKE) --no-print-directory testSuite

.PHONY: localApiTests-apiBugDemonstration-ocis
localApiTests-apiBugDemonstration-ocis: ## run apiBugDemonstration test suite with ocis storage
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=ocis \
	BEHAT_SUITE=apiBugDemonstration \
	$(MAKE) --no-print-directory testSuite

.PHONY: localApiTests-apiAccountsHashDifficulty-ocis
localApiTests-apiAccountsHashDifficulty-ocis: ## run apiAccountsHashDifficulty test suite with ocis storage
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=ocis \
	BEHAT_SUITE=apiAccountsHashDifficulty \
	$(MAKE) --no-print-directory testSuite

.PHONY: localApiTests-apiOcisSpecific-s3ng
localApiTests-apiOcisSpecific-s3ng: ## run apiOcisSPecific test suite with s3ng storage
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=s3ng \
	BEHAT_SUITE=apiOcisSpecific \
	START_CEPH=1 \
	$(MAKE) --no-print-directory testSuite


.PHONY: localApiTests-apiBasic-s3ng
localApiTests-apiBasic-s3ng: ## run apiBasic test suite with s3ng storage
	@TEST_SOURCE=ocis \
	STORAGE_DRIVER=s3ng \
	BEHAT_SUITE=apiBasic \
	START_CEPH=1 \
	$(MAKE) --no-print-directory testSuite

targets = $(addprefix Core-API-Tests-owncloud-storage-,$(PARTS))
.PHONY: $(targets)
$(targets):
	@$(eval RUN_PART=$(shell echo "$@" | tr -dc '0-9'))
	@TEST_SOURCE=oc10 \
	STORAGE_DRIVER=owncloud \
	RUN_PART=$(RUN_PART) \
	$(MAKE) --no-print-directory testSuite

targets = $(addprefix Core-API-Tests-ocis-storage-,$(PARTS))
.PHONY: $(targets)
$(targets):
	@$(eval RUN_PART=$(shell echo "$@" | tr -dc '0-9'))
	@TEST_SOURCE=oc10 \
	STORAGE_DRIVER=ocis \
	RUN_PART=$(RUN_PART) \
	$(MAKE) --no-print-directory testSuite

.PHONY: testSuite
testSuite: clean-docker-container
	@if [ -n "${START_CEPH}" ]; then \
		COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
		COMPOSE_FILE=src/ceph.yml \
		docker-compose run start_ceph; \
	fi; \
	COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
	COMPOSE_FILE=$(COMPOSE_FILE) \
	STORAGE_DRIVER=$(STORAGE_DRIVER) \
	TEST_SOURCE=$(TEST_SOURCE) \
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
	docker logs -f ocis-acceptance-tests_acceptance-tests_1 | less

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
	@$(MAKE) --no-print-directory -C ../. clean-tests

.PHONY: clean
clean: clean-docker-container clean-docker-volumes clean-files ## clean all
