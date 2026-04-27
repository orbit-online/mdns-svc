FROM golang:1.25.8-alpine
ARG TARGETARCH
ARG VERSION
RUN apk add pcsc-lite-libs pcsc-lite-dev gcc g++
COPY . /go/src/github.com/orbit-online/mdns-svc
WORKDIR /go/src/github.com/orbit-online/mdns-svc
RUN go get ./...
RUN go build -ldflags="-X main.VERSION=${VERSION} -s -linkmode external"

FROM scratch
COPY --from=0 /go/src/github.com/orbit-online/mdns-svc/mdns-svc /
