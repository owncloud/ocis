#!/bin/bash

set -e

CHART_REPO="$1"
if [[ -z "$CHART_REPO" ]]; then
    echo "[ERR] Chart directory argument missing. Usage: $0 <chart-repo-directory>"
    exit 1
fi

if [[ ! -d "$CHART_REPO" ]]; then
    echo "[ERR] Path not found: $CHART_REPO"
    exit 1
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$(cd $SCRIPT_DIR/../../.. && pwd)"

CFG_DIR="$ROOT/tests/config/k8s"
CHT_DIR="$CHART_REPO/charts/ocis"
TPL_DIR="$CHT_DIR/templates"

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
    sed -i '/virusscan:/{n;s|false|true|}' $CFG_DIR/values.yaml
fi

if [[ "$ENABLE_EMAIL" == "true" ]]; then
    sed -i '/emailNotifications:/{n;s|false|true|}' $CFG_DIR/values.yaml
fi

if [[ "$ENABLE_TIKA" == "true" ]]; then
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

if [[ "$ENABLE_AUTH_APP" == "true" ]]; then
    sed -i '/authapp:/{n;s|false|true|}' $CFG_DIR/values.yaml
fi

# [NOTE]
# Remove schema validation to add extra configs in values.yaml.
# Also this allows us to use fakeoffice as web-office server
rm "$CHT_DIR/values.schema.json"

# copy custom values file
cp $CFG_DIR/values.yaml "$CHT_DIR/ci/deployment-values.yaml"
