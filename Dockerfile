# Stage 1
FROM golang:1.21.0-alpine3.18 AS builder

WORKDIR /source
RUN apk --no-cache add git upx

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY src/ .
RUN go build -o ./dist/ ./...
RUN upx ./dist/*

# Stage 2
FROM alpine:3.18 AS runner

RUN adduser -D -h /opt/application -s /sbin/nologin runner
WORKDIR /opt/application
USER runner

COPY --from=builder /source/dist/ .
ENTRYPOINT ["/opt/application/cmd"]