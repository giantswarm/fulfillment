FROM gsoci.azurecr.io/giantswarm/alpine:3.24.1

RUN apk update && apk --no-cache add ca-certificates && update-ca-certificates

# architect/go-build emits one static binary per target platform
# (fulfillment-linux-amd64, fulfillment-linux-arm64) plus an unsuffixed copy of
# the linux/amd64 build. Copy the one matching buildx's TARGETARCH so the arm64
# image gets an arm64 binary. For a local build, produce it first:
#   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fulfillment-linux-amd64 .
ARG TARGETARCH
COPY ./fulfillment-linux-${TARGETARCH} /usr/local/bin/fulfillment/fulfillment
COPY ./content /usr/local/bin/fulfillment/content/

EXPOSE 8000
USER 9000:9000

WORKDIR /usr/local/bin/fulfillment
ENTRYPOINT ["/usr/local/bin/fulfillment/fulfillment"]
