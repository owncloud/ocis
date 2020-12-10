FROM webhippie/golang:1.15 as build

COPY ./ /ocis/
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk update && \
	apk upgrade --ignore musl-dev && \
	apk add make gcc bash && \
	rm -rf /var/cache/apk/*

WORKDIR /ocis/ocis
RUN make clean generate build


FROM alpine:3

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

RUN mkdir -p /var/lib/ocis/data
WORKDIR /var/lib/ocis

COPY --from=build /ocis/ocis/bin/ocis /usr/bin/ocis
