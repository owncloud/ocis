FROM owncloud/ocis:1.16.0 as ocis

# same base image as REVA: https://github.com/cs3org/reva/blob/master/Dockerfile.revad-eos#L39
# available imges: https://gitlab.cern.ch/dss/eos/container_registry/10190
FROM gitlab-registry.cern.ch/dss/eos/eos-all:4.8.57

COPY --from=ocis /usr/bin/ocis /usr/bin/ocis

COPY setup.sh /setup.sh
COPY entrypoint.sh /entrypoint.sh

RUN yum install -y \
    bc \
    initscripts \
    libmicrohttpd \
    protobuf3 \
    hostname \
    jemalloc \
    leveldb \
    nc \
    zeromq \
    redis \
    nss-pam-ldapd \
    nscd \
    authconfig \
    && yum clean all \
    && rm -rf /var/cache/yum

ENTRYPOINT [ "/entrypoint.sh" ]