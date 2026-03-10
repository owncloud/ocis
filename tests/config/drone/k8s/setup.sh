#!/bin/bash

set -e

ROOT="."
if [[ -n "$1" ]]; then
    ROOT="$1"
fi

CFG_DIR="$ROOT/tests/config/drone/k8s"
CHT_DIR="$ROOT/ocis-charts/charts/ocis"
TPL_DIR="$CHT_DIR/templates"

if [[ ! -d "$ROOT/ocis-charts" ]]; then
    echo "Error: ocis-charts not found in $ROOT/ocis-charts. Please clone it first."
    exit 1
fi

# patch ocis service templates
for service in "$TPL_DIR"/*/; do
    if [[ -f "$service/deployment.yaml" ]]; then
        if grep -qE 'ocis.caEnv' "$service/deployment.yaml"; then
            sed -i '/.*ocis.caEnv.*/a\{{- include "ocis.extraEnvs" . | nindent 12 }}' "$service/deployment.yaml"
            sed -i '/.*ocis.caPath.*/a\{{- include "ocis.extraVolMounts" . | nindent 12 }}' "$service/deployment.yaml"
            sed -i '/.*ocis.caVolume.*/a\{{- include "ocis.extraVolumes" . | nindent 8 }}' "$service/deployment.yaml"
        else
            sed -i '/env:/a\{{- include "ocis.extraEnvs" . | nindent 12 }}' "$service/deployment.yaml"
            sed -i '/volumeMounts:/a\{{- include "ocis.extraVolMounts" . | nindent 12 }}' "$service/deployment.yaml"
            sed -i '/volumes:/a\{{- include "ocis.extraVolumes" . | nindent 8 }}' "$service/deployment.yaml"
        fi
    fi
done

# copy custom template resources
cp -r $CFG_DIR/templates/* $TPL_DIR/

# add authbasic service
sed -i "/{{- define \"ocis.basicServiceTemplates\" -}}/a\  {{- \$_ := set .scope \"appNameAuthBasic\" \"authbasic\" -}}" $TPL_DIR/_common/_tplvalues.tpl

if [[ "$ENABLE_ANTIVIRUS" == "true" ]]; then
    # TODO: use external service
    cp -r $CFG_DIR/clamav $TPL_DIR/
    sed -i '/virusscan:/{n;s|false|true|}' $CFG_DIR/values.yaml
fi

if [[ "$ENABLE_EMAIL" == "true" ]]; then
    # TODO: use external service
    cp -r $CFG_DIR/mailpit $TPL_DIR/
    sed -i '/emailNotifications:/{n;s|false|true|}' $CFG_DIR/values.yaml
fi

if [[ "$ENABLE_TIKA" == "true" ]]; then
    # TODO: use external service
    cp -r $CFG_DIR/tika $TPL_DIR/
    sed -i 's|type: basic|type: tika|' $CFG_DIR/values.yaml
fi

if [[ "$ENABLE_WOPI" == "true" ]]; then
    sed -i '/appsIntegration:/{n;s|false|true|}' $CFG_DIR/values.yaml
    # patch collaboration service
    #  - allow dynamic wopi src
    sed -i -E "s|value: http://.*:9300|value: {{ \$officeSuite.wopiSrc }}|" $TPL_DIR/collaboration/deployment.yaml
fi

if [[ "$ENABLE_OCM" == "true" ]]; then
    sed -i '/ocm:/{n;s|false|true|}' $CFG_DIR/values.yaml
fi

# move custom values file
cp $CFG_DIR/values.yaml "$CHT_DIR/ci/deployment-values.yaml"
