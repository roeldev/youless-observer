# golang version 1 = latest 1.x.x
ARG GOLANG_VERSION="1"

###############################################################################
FROM golang:${GOLANG_VERSION} as builder

ENV CGO_ENABLED=0  \
    GOOS=linux  \
    GOARCH=amd64

COPY ./ /go/build
WORKDIR /go/build/

ARG APP=observer
RUN set -eux  \
    && mkdir -p /root-out/  \
    && cp cmd/${APP}/.env /root-out/.env

RUN set -eux  \
    && go mod download  \
    && go test  \
    && go build  \
        -tags=notrace  \
        -ldflags="-s -w"  \
        -o="/root-out/main"  \
        ./cmd/${APP}/...

###############################################################################
# create actual image
FROM gcr.io/distroless/static
COPY --from=builder /root-out/ /

EXPOSE 2512
USER 1000
ENTRYPOINT ["/main"]
