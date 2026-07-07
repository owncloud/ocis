FROM owncloudops/nginx@sha256:8ca0c90c2f537fef6195ba8f5ebe6d6c00fe0a0c545f25ce8acd418a831a19a4

LABEL maintainer="ownCloud GmbH <devops@owncloud.com>" \
  org.opencontainers.image.title="ownCloud Web" \
  org.opencontainers.image.vendor="ownCloud GmbH" \
  org.opencontainers.image.authors="ownCloud GmbH" \
  org.opencontainers.image.description="ownCloud Web - User interface for ownCloud" \
  org.opencontainers.image.licenses="Apache-2.0" \
  org.opencontainers.image.documentation="https://github.com/owncloud/ocis" \
  org.opencontainers.image.url="https://hub.docker.com/r/owncloud/web" \
  org.opencontainers.image.source="https://github.com/owncloud/ocis"

RUN rm -f /var/lib/nginx/html/*

ADD docker/overlay /
ADD dist/ /var/lib/nginx/html

EXPOSE 8080

USER nginx

STOPSIGNAL SIGTERM

CMD ["nginx", "-g", "daemon off;"]
WORKDIR /var/lib/nginx/html
