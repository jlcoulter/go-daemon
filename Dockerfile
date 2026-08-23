# Build stage
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /daemon ./cmd/daemon

# Final stage
FROM scratch
COPY --from=builder /daemon /daemon
EXPOSE 2222 8080
ENTRYPOINT ["/daemon"]