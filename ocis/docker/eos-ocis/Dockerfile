FROM centos:7 as build
ARG BRANCH=master

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
RUN git clone https://github.com/owncloud/ocis.git
RUN cd ocis && git checkout ${BRANCH}
RUN cd ocis && git log -n 1
RUN cd ocis/ocis && go build -v -o ../bin/ocis -toolexec '/usr/bin/time -v' ./cmd/ocis
# RUN cd ocis && make generate bin/ocis

# build user lookup test program
RUN mkdir -p /build
COPY user.go /build/user.go
RUN cd /build && go build -o user user.go


FROM owncloud/eos-base:4.8.26
RUN mkdir -p /usr/local/bin
COPY --from=build ocis/bin/ocis /usr/local/bin/
RUN chmod +x /usr/local/bin/ocis

COPY --from=build /build/user /usr/local/bin/user
RUN chmod +x /usr/local/bin/user

COPY setup /setup
COPY entrypoint /entrypoint
