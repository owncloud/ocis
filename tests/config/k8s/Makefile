
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
		-v $(LOGS_DIR):/logs \
		-v $(abspath ./traefik-config.yaml):/var/lib/rancher/k3s/server/manifests/traefik-config.yaml

.PHONY: prepare-charts
prepare-charts:
	@test -d "$(CHART_REPO)" || git clone --single-branch -b main --depth 1 https://github.com/owncloud/ocis-charts.git $(CHART_REPO)
	bash setup.sh $(CHART_REPO)
	kubectl create namespace $(DOMAIN) || true
	kubectl create configmap -n $(DOMAIN) ocis-fonts-ttf --from-file=../ci/NotoSans.ttf || true
	kubectl create secret generic -n $(DOMAIN) ocis-admin-user \
		--from-literal=user-id=3c5f1c9e-8a4b-4d2e-9b7a-6f0e1d2c3b4a --from-literal=password=admin || true
	sed 's|https://ocis-server:9200|https://$(DOMAIN)|g' ../ci/ocis-config.json > $(CHART_REPO)/web-config.json
	@if [ "$(ENABLE_WOPI)" = "true" ]; then sed -i 's|"files",|"files", "external",|' $(CHART_REPO)/web-config.json; fi
	@if [ "$(ENABLE_OCM)" = "true" ]; then sed -i 's|"files",|"files", "ocm",|' $(CHART_REPO)/web-config.json; fi
	kubectl create configmap -n $(DOMAIN) ocis-web-config \
		--from-file=web-config.json=$(CHART_REPO)/web-config.json || true
	kubectl create configmap -n $(DOMAIN) ocis-translations \
		--from-file=de_LC_MESSAGES_notifications.po=../translations/de/LC_MESSAGES/notifications.po \
		--from-file=es_LC_MESSAGES_notifications.po=../translations/es/LC_MESSAGES/notifications.po || true

.PHONY: deploy-ocis
deploy-ocis:
	@sed -i 's|externalDomain:.*|externalDomain: $(DOMAIN)|' $(CHART_REPO)/charts/ocis/ci/deployment-values.yaml
	helm install -n $(DOMAIN) --create-namespace \
		--values $(CHART_REPO)/charts/ocis/ci/deployment-values.yaml \
		--wait --timeout 5m0s ocis $(CHART_REPO)/charts/ocis/

.PHONY: show-logs
show-logs:
	@tail -f $(LOGS_DIR)/*.log

.PHONY: cleanup
cleanup:
	k3d cluster delete local-ocis
	rm -rf $(LOGS_DIR)
	rm -rf $(CHART_REPO)
