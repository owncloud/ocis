FROM owncloud/eos-base:4.8.26
RUN mkdir -p /usr/local/bin
COPY dist/binaries/ocis-linux-amd64 /usr/local/bin/ocis
RUN chmod +x /usr/local/bin/ocis

COPY docker/eos-ocis/setup /setup
COPY docker/eos-ocis/entrypoint /entrypoint
