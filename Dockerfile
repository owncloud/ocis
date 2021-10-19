FROM owncloudci/golang:1.17 as build

RUN apk update && \
        apk add --update npm

RUN npm install --global yarn

COPY ./ /ocis/

WORKDIR /ocis/ocis
RUN make clean generate build


FROM alpine:3.13

RUN apk update && \
	apk upgrade && \
	apk add ca-certificates mailcap && \
        apk add --update npm && \
	rm -rf /var/cache/apk/* && \
	echo 'hosts: files dns' >| /etc/nsswitch.conf

RUN npm install --global yarn

LABEL maintainer="ownCloud GmbH <devops@owncloud.com>" \
  org.label-schema.name="ownCloud Infinite Scale" \
  org.label-schema.vendor="ownCloud GmbH" \
  org.label-schema.schema-version="1.0"

ENTRYPOINT ["/usr/bin/ocis"]
CMD ["server"]

COPY --from=build /ocis/ocis/bin/ocis /usr/bin/ocis
