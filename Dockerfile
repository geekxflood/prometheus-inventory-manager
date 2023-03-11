FROM golang:1.20-alpine AS builder

ARG ARCH=amd64

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH
ENV GO_VERSION 1.20
ENV GO111MODULE on
ENV CGO_ENABLED=0

# Build dependencies
WORKDIR /go/src/
COPY . .
RUN apk update && apk add make git
RUN go get ./...
RUN mkdir /go/src/build 
RUN go build -a -gcflags=all="-l -B" -ldflags="-w -s" -o build/prim ./...

# Second stage
FROM alpine:3.17

RUN apk update

# Copy your custom certificate file to the container
# COPY my-cert.crt /usr/local/share/ca-certificates/my-cert.crt

# Install the ca-certificates package
# RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /go/src/build/prim /usr/local/bin/prim
CMD ["/usr/local/bin/prim"]

