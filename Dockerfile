FROM golang:latest as builder

WORKDIR /app 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" ./cmd/resume-pdf-converter

FROM browserless/chrome:1-chrome-stable

WORKDIR /app

COPY --from=builder /app/resume-pdf-converter /usr/bin/

ENTRYPOINT ["resume-pdf-converter"]
