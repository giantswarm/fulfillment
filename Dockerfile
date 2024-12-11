FROM gsoci.azurecr.io/giantswarm/alpine:3.21.0

RUN apk update && apk --no-cache add ca-certificates && update-ca-certificates

RUN mkdir -p /usr/loca/bin/fulfillment/
RUN mkdir -p /usr/loca/bin/fulfillment/content/

ADD ./fulfillment /usr/local/bin/fulfillment/fulfillment
ADD ./content /usr/local/bin/fulfillment/content/

EXPOSE 8000
USER 9000:9000

WORKDIR /usr/local/bin/fulfillment
ENTRYPOINT ["/usr/local/bin/fulfillment/fulfillment"]
