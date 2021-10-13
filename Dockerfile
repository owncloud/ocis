FROM owncloudci/golang:1.17 as build

COPY ./ /ocis/

WORKDIR /ocis/ocis
RUN make clean generate build


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
