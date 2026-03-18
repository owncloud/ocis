
LOGS_DIR := $(abspath ./logs)
CHART_REPO := $(abspath ./ocis-charts)
HOST_IP := $(shell ip route get 1 | awk '{print $$7}')

OCM ?= false
DOMAIN := ocis-server
ifeq ($(OCM), true)
	DOMAIN := federation-ocis-server
endif

.PHONY: create-cluster
create-cluster:
	@echo $(LOGS_DIR)
	@mkdir -p $(LOGS_DIR)
	@chmod 777 $(LOGS_DIR)
	k3d cluster create local-ocis \
		-p '80:80@loadbalancer' -p '443:443@loadbalancer' -p '9100-9399:30100-30399@loadbalancer' \
		--k3s-arg '--tls-san=k3d@server:*' --k3s-arg '--disable=metrics-server@server:*' \
		--host-alias $(HOST_IP):ocis-server --host-alias $(HOST_IP):federation-ocis-server \
		-v $(LOGS_DIR):/logs

.PHONY: prepare-charts
prepare-charts:
	@test -d "$(CHART_REPO)" || git clone --single-branch -b main --depth 1 https://github.com/owncloud/ocis-charts.git $(CHART_REPO)
	bash setup.sh $(CHART_REPO)
	kubectl create namespace $(DOMAIN) || true
	kubectl create configmap -n $(DOMAIN) ocis-fonts-ttf --from-file=../drone/NotoSans.ttf || true

.PHONY: deploy-ocis
deploy-ocis:
	@sed -i 's|externalDomain:.*|externalDomain: $(DOMAIN)|' $(CHART_REPO)/charts/ocis/ci/deployment-values.yaml
	helm install -n $(DOMAIN) --create-namespace \
		--values $(CHART_REPO)/charts/ocis/ci/deployment-values.yaml \
		--wait --timeout 5m0s ocis $(CHART_REPO)/charts/ocis/

cleanup:
	k3d cluster delete local-ocis
	rm -rf $(LOGS_DIR)
	rm -rf $(CHART_REPO)
