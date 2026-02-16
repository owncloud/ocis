# custom Dockerfile required to run ociswrapper command
# mounting 'ociswrapper' binary doesn't work with image 'amd64/alpine:3.17' (busybox based)

ARG OCIS_IMAGE_TAG
FROM owncloud/ocis:${OCIS_IMAGE_TAG} as ocis

FROM ubuntu:22.04
COPY --from=ocis /usr/bin/ocis /usr/bin/ocis

COPY ["./serve-ocis.sh", "/usr/bin/serve-ocis"]
RUN chmod +x /usr/bin/serve-ocis

ENTRYPOINT [ "serve-ocis" ]