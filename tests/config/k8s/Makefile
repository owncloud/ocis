
LOGS_DIR := $(abspath ./logs)
CHART_REPO := $(abspath ./ocis-charts)

DOMAIN ?= ocis-server

.PHONY: create-cluster
create-cluster:
	@echo $(LOGS_DIR)
	@mkdir -p $(LOGS_DIR)
	k3d cluster create local-ocis \
	-p '80:80@loadbalancer' -p '443:443@loadbalancer' -p '9100-9399:30100-30399@loadbalancer' \
	--k3s-arg '--tls-san=k3d@server:*' --k3s-arg '--disable=metrics-server@server:*' \
	-v $(LOGS_DIR):/logs

.PHONY: setup-charts
setup-charts:
	@[ ! -d "$(CHART_REPO)" ] && git clone --single-branch -b main --depth 1 https://github.com/owncloud/ocis-charts.git $(CHART_REPO)
	bash setup.sh $(CHART_REPO)
	kubectl create namespace ocis || true
	kubectl create configmap -n ocis ocis-fonts-ttf --from-file=../drone/NotoSans.ttf || true

.PHONY: deploy-ocis
deploy-ocis:
	@cd $(CHART_REPO)
	@sed -i 's|externalDomain:.*|externalDomain: $(DOMAIN)|' $(CHART_REPO)/charts/ocis/ci/deployment-values.yaml
	make helm-install-atomic

cleanup:
	k3d cluster delete local-ocis
	@rm -rf $(LOGS_DIR)
