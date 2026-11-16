# Stage 1: Modules caching
FROM golang:1.22 AS modules
WORKDIR /modules
COPY go.mod go.sum ./
RUN go mod download

# Stage 2: Build
FROM golang:1.22 AS builder
COPY --from=modules /go/pkg /go/pkg
WORKDIR /src/rms-media-discovery
COPY . .
# Build playwright CLI with the version pinned in go.mod (used later to install browsers)
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/playwright github.com/mxschmitt/playwright-go/cmd/playwright
# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=`git tag --sort=-version:refname | head -n 1`" -o rms-media-discovery -a -installsuffix cgo rms-media-discovery.go

# Stage 3: Final
FROM ubuntu:noble
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/bin/playwright /usr/local/bin/playwright
# Install the playwright driver, system dependencies and headless chromium
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates tzdata \
 && playwright install --with-deps chromium \
 && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/rms-media-discovery/rms-media-discovery .
COPY --from=builder /src/rms-media-discovery/configs/rms-media-discovery.json /etc/rms/
EXPOSE 8080/tcp
EXPOSE 2112/tcp
CMD ["./rms-media-discovery"]
