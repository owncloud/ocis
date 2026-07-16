#!/usr/bin/env bash

set -euo pipefail

NAMESPACE="${NAMESPACE:-ocis-server}"

declare -A DEBUG_PORTS=(
    [activitylog]=9197
    [antivirus]=9277
    [appregistry]=9243
    [audit]=9229
    [authapp]=9245
    [authbasic]=9147
    [authmachine]=9167
    [authservice]=9198
    [clientlog]=9260
    [collaboration]=9304
    [eventhistory]=9270
    [frontend]=9141
    [gateway]=9143
    [graph]=9124
    [groups]=9161
    [idm]=9239
    [idp]=9134
    [nats]=9234
    [notifications]=9174
    [ocdav]=9163
    [ocm]=9281
    [ocs]=9114
    [postprocessing]=9255
    [proxy]=9205
    [search]=9224
    [settings]=9194
    [sharing]=9151
    [sse]=9139
    [storagepubliclink]=9179
    [storageshares]=9156
    [storagesystem]=9217
    [storageusers]=9159
    [thumbnails]=9189
    [userlog]=9214
    [users]=9145
    [web]=9104
    [webdav]=9119
    [webfinger]=9279
)

declare -A GRPC_PORTS=(
    [appregistry]=9242
    [authapp]=9246
    [authbasic]=9146
    [authmachine]=9166
    [authservice]=9616
    [collaboration]=9301
    [eventhistory]=8080
    [gateway]=9142
    [groups]=9160
    [ocm]=9282
    [search]=9220
    [settings]=9191
    [sharing]=9150
    [storagepubliclink]=9178
    [storageshares]=9154
    [storagesystem]=9215
    [storageusers]=9157
    [thumbnails]=9185
    [users]=9144
)

should_expose() {
    case "$1" in
        antivirus)
            [[ "${ENABLE_ANTIVIRUS:-false}" == "true" ]]
            ;;
        collaboration)
            [[ "${ENABLE_WOPI:-false}" == "true" ]]
            ;;
		collaboration)
            [[ "${ENABLE_WOPI:-false}" == "true" ]]
            ;;
        ocm)
            [[ "${ENABLE_OCM:-false}" == "true" ]]
            ;;
        authapp)
            [[ "${ENABLE_AUTH_APP:-false}" == "true" ]]
            ;;
        notifications)
            [[ "${ENABLE_EMAIL:-false}" == "true" ]]
            ;;
        search)
            [[ "${ENABLE_TIKA:-false}" == "true" ]]
            ;;
        *)
            return 0
            ;;
    esac
}

expose() {
    local deployment=$1
    local service=$2
    local port=$3

    if ! should_expose "$deployment"; then
        echo "[SKIP] $deployment disabled"
        return
    fi

    if kubectl -n "$NAMESPACE" get svc "$service" >/dev/null 2>&1; then
        echo "[SKIP] $service already exists"
        return
    fi

    kubectl -n "$NAMESPACE" expose deployment "$deployment" \
        --name="$service" \
        --port="$port" \
        --target-port="$port"
}

for svc in "${!DEBUG_PORTS[@]}"; do
    expose "$svc" "${svc}-debug" "${DEBUG_PORTS[$svc]}"
done

for svc in "${!GRPC_PORTS[@]}"; do
    expose "$svc" "${svc}-grpc" "${GRPC_PORTS[$svc]}"
done