# Build on the native build platform and cross-compile to the target arch so
# multi-arch builds don't pay for emulating the Go toolchain.
FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/go:latest AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" ./cmd/resume-pdf-converter

FROM cgr.dev/chainguard/wolfi-base:latest

WORKDIR /app

# Install Chromium from the Wolfi package repo. Lands at /usr/bin/chromium,
# which is on PATH so chromedp can discover it directly.
RUN apk add --no-cache chromium

COPY --from=builder /app/resume-pdf-converter /usr/bin/

ENV CHROME_PATH=/usr/bin/chromium

# Run as a non-root user; Chrome runs with --no-sandbox so no SYS_ADMIN needed.
USER 65532

ENTRYPOINT ["resume-pdf-converter"]
