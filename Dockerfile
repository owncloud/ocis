# Please consider to not use this Dockerfile.
# If you want to build an image from source,
# you might want to run following command instead:
# `make -C ocis dev-docker`
# It will build a `owncloud/ocis:dev` image for you.

# If you still want to build oCIS using this Dockerfile
# you can do it by running following command:
# `docker build -t owncloud/ocis:custom .`

FROM owncloudci/nodejs:14 as generate

COPY ./ /ocis/

WORKDIR /ocis/ocis
RUN make ci-node-generate

FROM owncloudci/golang:1.17 as build

COPY --from=generate /ocis /ocis

WORKDIR /ocis/ocis
RUN make ci-go-generate build

FROM alpine:3.13

RUN apk update && \
	apk upgrade && \
	apk add ca-certificates mailcap && \
	rm -rf /var/cache/apk/* && \
	echo 'hosts: files dns' >| /etc/nsswitch.conf

LABEL maintainer="ownCloud GmbH <devops@owncloud.com>" \
	org.label-schema.name="ownCloud Infinite Scale" \
	org.label-schema.vendor="ownCloud GmbH" \
	org.label-schema.schema-version="1.0"

ENTRYPOINT ["/usr/bin/ocis"]
CMD ["server"]

COPY --from=build /ocis/ocis/bin/ocis /usr/bin/ocis
