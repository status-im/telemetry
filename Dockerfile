FROM golang:1.22-alpine3.19 AS builder

# Create a non-root user
RUN addgroup -g 1001 -S iamgroup && \
    adduser -u 1001 -S telemetry -G iamgroup

WORKDIR /go/src/github.com/status-im/telemetry

# Copy go mod and sum files for caching
COPY go.mod go.sum ./

# Using mounts to cache dependencies
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download -x

COPY . .

# Build the binary with static linking and stripped symbols for smaller size
RUN GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags "-s -w -extldflags -static" \
    -o build/server cmd/server/main.go

# Copy the binary to final image
FROM alpine:3.19

LABEL maintainer="jakub@status.im" \
    source="https://github.com/status-im/telemetry" \
    description="Opt-in message reliability metrics service" \
    commit="unknown"

# Copy the /etc/passwd file from the builder stage to provide non-root user information
COPY --from=builder /etc/passwd /etc/passwd

# Copy the compiled application binary from the build stage to the final image
COPY --from=builder /go/src/github.com/status-im/telemetry/build/server /usr/local/bin/telemetry

USER telemetry

EXPOSE 8080/tcp

ENTRYPOINT ["/usr/local/bin/telemetry"]
CMD ["-help"]
