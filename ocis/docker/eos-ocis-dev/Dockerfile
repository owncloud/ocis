FROM owncloud/eos-base:4.8.26

ENV CGO_ENABLED=1
ENV GOOS=linux

RUN rpm --rebuilddb && yum -y install \
    wget \
    time \
    make \
    gcc \
    git

RUN wget -q https://dl.google.com/go/go1.15.7.linux-amd64.tar.gz
RUN mkdir -p /usr/local/bin
RUN tar xf go1.15.7.linux-amd64.tar.gz -C /usr/local
RUN ln -s /usr/local/go/bin/* /usr/local/bin

COPY entrypoint /entrypoint
COPY start-ldap /start-ldap

VOLUME [ "/ocis" ]
WORKDIR /ocis


