# Stage 1: Build the binary (Bumped from 1.23 to 1.26 to satisfy go.mod requirements)
FROM owncloudci/golang:1.25 AS builder

RUN apk add --no-cache make git

WORKDIR /ociswrapper

# Copy everything from the current local folder (tests/ociswrapper) into the container
COPY . .

# Run the local make target to build the binary
RUN make build

# Stage 2: Minimal runtime environment
FROM alpine:3.23.4

WORKDIR /ociswrapper

# Copy only the compiled binary
COPY --from=builder /ociswrapper/bin/ociswrapper /usr/local/bin/ociswrapper

EXPOSE 5200

ENTRYPOINT ["ociswrapper"]