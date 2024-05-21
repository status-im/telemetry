FROM golang:1.15-alpine AS builder

RUN apk add --no-cache make

RUN mkdir -p /go/src/github.com/status-im/telemetry
WORKDIR /go/src/github.com/status-im/telemetry
ADD . .
RUN make build

# Copy the binary to the second image
FROM alpine:3.14

LABEL maintainer="jakub@status.im"
LABEL source="https://github.com/status-im/telemetry"
LABEL description="Opt-in message reliability metrics service"
LABEL commit="unknown"

COPY --from=builder /go/src/github.com/status-im/telemetry/build/server /usr/local/bin/telemetry

EXPOSE 8080/tcp

ENTRYPOINT ["/usr/local/bin/telemetry"]
CMD ["-help"]
